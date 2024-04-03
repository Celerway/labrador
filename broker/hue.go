package broker

import (
	"fmt"
	"github.com/celerway/labrador/gohue"
	"log/slog"
	"os"
)

func newHueClient(logger *slog.Logger) (*gohue.HueClient, error) {
	bridge, ok := os.LookupEnv("HUE_BRIDGE")
	if !ok {
		return nil, fmt.Errorf("HUE_BRIDGE not set")
	}
	username, ok := os.LookupEnv("HUE_USERNAME")
	if !ok {
		return nil, fmt.Errorf("HUE_USERNAME not set")
	}
	c, err := gohue.New(bridge, username, logger)
	if err != nil {
		return nil, fmt.Errorf("gohue.NewClient: %w", err)
	}
	return c, nil
}
