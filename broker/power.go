package broker

import (
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

	switch action {
	case "control":
		err := s.onPowerControl(device, pk.Payload)
		if err != nil {
			s.logger.Warn("onPowerControl", "device", device, "error", err)
		}
	case "status":
		err := s.onPowerStatus(device, pk.Payload)
		if err != nil {
			s.logger.Warn("onPowerStatus", "device", device, "error", err)
		}
	default:
		s.logger.Warn("unknown power action", "action", action, "device", device, "topic", pk.TopicName)
	}
}

func (s *State) onPowerControl(device string, payload []byte) error {
	// payload is the desired state of the device.
	// This is where you would send the control command to the device.
	return nil
}

func (s *State) onPowerStatus(device string, payload []byte) error {
	// payload is the current state of the device.
	// This is where you would update the state of the device in the broker.
	return nil
}
