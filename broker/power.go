package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/celerway/labrador/msgs"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"regexp"
)

type PowerDevice struct {
	State  bool // true = on, false = off
	Father *State
}

type PowerStatus struct {
}

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
		dev, ok := s.internalPDs[device]
		if !ok {
			s.logger.Info("ignoring non-builtin device", "device", device, "topic", pk.TopicName)
			return
		}
		err := dev.onPowerControl(device, pk.Payload)
		if err != nil {
			s.logger.Warn("onPowerControl", "device", device, "error", err)
		}
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

func (s *State) NewPowerDevice(deviceID string) error {
	if _, ok := s.internalPDs[deviceID]; ok {
		return errors.New("device already exists")
	}
	pd := PowerDevice{
		State:  false,
		Father: s,
	}
	s.internalPDs[deviceID] = &pd
	return nil
}

func (pd *PowerDevice) onPowerControl(device string, payloadBytes []byte) error {
	huec := pd.Father.bridgeConn
	if huec == nil {
		return errors.New("no hue client")
	}
	var payload msgs.PowerControl
	err := json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}
