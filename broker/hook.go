package broker

import (
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"sync"
)

/*
The monitor hook extracts information from the broker so the web dashboard can display it.
*/

type MonitorHook struct {
	mu sync.Mutex
	mqtt.HookBase
	config    *MonitorHookOptions
	clientMap map[string]*mqtt.Client
	msgs      *circularBuffer
}

type MonitorHookOptions struct {
	Server *mqtt.Server
}

func (h *MonitorHook) ID() string {
	return "monitor"
}

func (h *MonitorHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnPublished,
	}, []byte{b})
}

func (h *MonitorHook) Init(config any) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := config.(*MonitorHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}
	h.config = config.(*MonitorHookOptions)
	if h.config.Server == nil {
		return mqtt.ErrInvalidConfigType
	}
	h.clientMap = make(map[string]*mqtt.Client)
	h.msgs = newBuffer(10)
	return nil
}

func (h *MonitorHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	h.Log.Info("client connected", "client", cl.ID)
	h.mu.Lock()
	h.clientMap[cl.ID] = cl
	h.mu.Unlock()
	return nil
}

func (h *MonitorHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info("client disconnected", "client", cl.ID)
	h.mu.Lock()
	delete(h.clientMap, cl.ID)
	h.mu.Unlock()
}

func (h *MonitorHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info("packet published", "client", cl.ID)
	h.msgs.push(pk)
}

func (h *MonitorHook) Close() error {
	h.Log.Info("MonitorHook.Close")
	return nil
}

func (h *MonitorHook) clients() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	cs := make([]string, 0, len(h.clientMap))
	for id := range h.clientMap {
		cs = append(cs, id)
	}
	return cs
}

func (h *MonitorHook) messages() []packets.Packet {
	return h.msgs.get()
}
