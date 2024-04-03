package broker

import (
	"context"
	"encoding/json"
	"github.com/celerway/labrador/msgs"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"regexp"
)

var powerTopicRx = regexp.MustCompile(`lab/power/(\w+)/(control|status)`)

func (s *State) onPower(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	// payload in pk.Payload

	// Extract the device and action from the topic.
	m := powerTopicRx.FindStringSubmatch(pk.TopicName)
	if m == nil {
		s.logger.Warn("parsing power topic failed", "topic", pk.TopicName, "client id", cl.ID)
		return
	}
	device := m[1]
	action := m[2]
	s.logger.Debug("onPower", "device", device, "action", action, "topic", pk.TopicName)
	switch action {
	case "control":
		var dev msgs.PowerControl
		err := json.Unmarshal(pk.Payload, &dev)
		if err != nil {
			s.logger.Warn("json.Unmarshal", "error", err, "topic", pk.TopicName)
			return
		}
		for _, pd := range s.HueBridge.GetPlugs() {
			if pd == device {
				err := s.HueBridge.SetPlug(context.TODO(), device, dev.Power)
				if err != nil {
					s.logger.Warn("onPowerControl", "device", device, "error", err)
				}
				return
			}
		}
		s.logger.Warn("device not found", "device", device, "action", action, "topic", pk.TopicName)
	case "status":
		var dev msgs.PowerStatus
		err := json.Unmarshal(pk.Payload, &dev)
		if err != nil {
			s.logger.Warn("json.Unmarshal", "error", err)
			return
		}
		s.logger.Debug("onPowerStatus", "device", device, "status", dev.Power, "error", dev.Error)
		s.pdStatus[device] = dev
	default:
		s.logger.Warn("unknown power action", "action", action, "device", device, "topic", pk.TopicName)
	}
}
