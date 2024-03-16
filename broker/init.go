package broker

import (
	"context"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"log/slog"
	"sync"
)

type State struct {
	mu      sync.Mutex
	logger  *slog.Logger
	server  *mqtt.Server
	clients map[string]*mqtt.Client
}

func New(logger *slog.Logger) *State {
	server := mqtt.New(&mqtt.Options{
		Logger: logger,
	})

	s := &State{
		logger: logger,
		server: server,
	}
	return s
}

func (s *State) Run(ctx context.Context) error {
	// Allow all connections.
	err := s.server.AddHook(new(auth.AllowHook), nil)
	if err != nil {
		return fmt.Errorf("mqtt server.AddHook: %w", err)
	}
	err = s.server.AddHook(new(MonitorHook), &MonitorHookOptions{
		Server: s.server,
	})
	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", nil)
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
	// Run server until context is cancelled or an error occurs.
	select {
	case <-ctx.Done():
		return s.server.Close()
	case err := <-failCh:
		return fmt.Errorf("mqtt server.Serve: %w", err)
	}
}
