package broker

import (
	"fmt"
	"github.com/Celerway/labrador/gohue"
)

func newHueClient(server string) (*gohue.Client, error) {
	c, err := gohue.NewClient(server)
	if err != nil {
		return nil, fmt.Errorf("gohue.NewClient: %w", err)
	}
	return c, nil
}
