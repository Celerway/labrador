package broker

import (
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type MonitorHook struct {
	mqtt.HookBase
	config *MonitorHookOptions
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
	return nil
}

func (h *MonitorHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info("client disconnected", "client", cl.ID)
}

func (h *MonitorHook) Close() error {
	h.Log.Info("MonitorHook.Close")
	return nil
}
