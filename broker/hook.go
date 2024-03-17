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
	config  *MonitorHookOptions
	clients map[string]*mqtt.Client
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
	}, []byte{b})
}

func (h *MonitorHook) Init(config any) error {
	h.Log.Info("MonitorHook.Init", "config", config)
	if _, ok := config.(*MonitorHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}
	h.config = config.(*MonitorHookOptions)
	if h.config.Server == nil {
		return mqtt.ErrInvalidConfigType
	}
	return nil
}

func (h *MonitorHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	h.Log.Info("client connected", "client", cl.ID)
	h.mu.Lock()
	h.clients[cl.ID] = cl
	h.mu.Unlock()
	return nil
}

func (h *MonitorHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info("client disconnected", "client", cl.ID)
	h.mu.Lock()
	delete(h.clients, cl.ID)
	h.mu.Unlock()
}

func (h *MonitorHook) Close() error {
	h.Log.Info("MonitorHook.Close")
	return nil
}
