package broker

import (
	"context"
	"fmt"
	"github.com/celerway/labrador/gohue"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"log/slog"
	"sync"
)

type State struct {
	mu         sync.Mutex
	logger     *slog.Logger
	server     *mqtt.Server
	bridgeConn *gohue.Client
	monitor    *MonitorHook
}

type PowerDevice struct {
	DeviceID string
	State    bool // true = on, false = off
	Internal bool
}

const (
	powerSub int = iota + 1
	storageSub
)

func New(logger *slog.Logger) *State {
	monitor := &MonitorHook{
		clients: make(map[string]*mqtt.Client),
	}

	server := mqtt.New(&mqtt.Options{
		Logger:       logger,
		InlineClient: true, // enable inline client support, allows us to directly publish and subscribe
	})

	s := &State{
		logger:  logger,
		server:  server,
		monitor: monitor,
	}
	return s
}

func (s *State) Run(ctx context.Context) error {
	// Allow all connections.
	err := s.server.AddHook(new(auth.AllowHook), nil)
	if err != nil {
		return fmt.Errorf("mqtt server.AddHook: %w", err)
	}
	err = s.server.AddHook(s.monitor, &MonitorHookOptions{Server: s.server})
	if err != nil {
		return fmt.Errorf("mqtt server.AddHook: %w", err)
	}
	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("tcp", ":1883", nil)
	err = s.server.AddListener(tcp)
	if err != nil {
		return fmt.Errorf("mqtt server.AddListener: %w", err)
	}
	failCh := make(chan error)
	go func() {
		err := s.server.Serve()
		if err != nil {
			failCh <- err
		}
	}()
	err = s.server.Subscribe("lab/power/#", powerSub, s.onPower)
	if err != nil {
		return fmt.Errorf("mqtt server.Subscribe(lab/power/#): %w", err)
	}
	err = s.server.Subscribe("lab/storage/#", storageSub, s.onStorage)
	if err != nil {
		return fmt.Errorf("mqtt server.Subscribe(lab/storage/#): %w", err)
	}
	// set up the hue client
	s.bridgeConn, err = newHueClient("http://1.2.3.4:80/")
	if err != nil {
		return fmt.Errorf("newHueClient: %w", err)
	}

	// Run server until context is cancelled or an error occurs.
	select {
	case <-ctx.Done():
		return s.server.Close()
	case err := <-failCh:
		return fmt.Errorf("mqtt server.Serve: %w", err)
	}
}

func (s *State) onStorage(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	fmt.Println("onStorage")
}
