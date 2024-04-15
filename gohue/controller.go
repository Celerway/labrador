package gohue

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"github.com/celerway/labrador/gohue/color"
	"log/slog"
	"net"
	"net/http"
	"sort"
)

type HueClient struct {
	Client *Client
	plugs  map[string]Plug
	logger *slog.Logger
	strips []Strip
}

type Plug struct {
	knownState bool
	id         string
}

type Strip struct {
	name string
	id   string
}

//go:embed hue_bridge_cacert.pem
var bridgeCaCert []byte

func newHttpClient() *http.Client {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(bridgeCaCert)
	// Set up HTTPS client
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return httpClient
}

func New(server, username string, logger *slog.Logger) (*HueClient, error) {
	ip, err := resolveIP(server)
	if err != nil {
		return nil, fmt.Errorf("resolveIP: %w", err)
	}
	baseUrl := fmt.Sprintf("https://%s/", ip)
	c, err := NewClient(baseUrl, WithHTTPClient(newHttpClient()), WithRequestEditorFn(makeEditorFn(username)))
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	return &HueClient{
		Client: c,
		logger: logger,
	}, nil
}

// resolveIP resolves the IP address of the Hue bridge, given the hostname.
func resolveIP(server string) (string, error) {
	ips, err := net.LookupIP(server)
	if err != nil {
		return "", fmt.Errorf("net.LookupIP: %w", err)
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found")
	}
	// return the first IPv4 address we find:
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("no IPv4 addresses found")
}

func makeEditorFn(username string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("hue-application-key", username)
		return nil
	}
}

func (c *HueClient) Load(ctx context.Context) error {
	resp, err := c.Client.GetLights(ctx)
	if err != nil {
		return fmt.Errorf("c.Client.GetLights: %w", err)
	}
	parsed, err := ParseGetLightsResponse(resp)
	if err != nil {
		return fmt.Errorf("ParseGetLightsResponse: %w", err)
	}
	if parsed.JSON200.Data == nil {
		return fmt.Errorf("no lights found")
	}
	lights := *parsed.JSON200.Data
	if len(lights) == 0 {
		return fmt.Errorf("no lights found")
	}
	c.plugs = make(map[string]Plug)
	for _, light := range lights {
		if light.Id == nil {
			return fmt.Errorf("light id is nil")
		}
		if light.Metadata.Name == nil {
			return fmt.Errorf("light name is nil")
		}
		if light.Metadata.Archetype == nil {
			return fmt.Errorf("light archetype is nil")
		}

		switch *light.Metadata.Archetype {
		case "plug":
			c.plugs[*light.Metadata.Name] = Plug{
				knownState: *light.On.On,
				id:         *light.Id,
			}
			c.logger.Info("found light", "name", *light.Metadata.Name, "id", *light.Id)
		case "hue_lightstrip":
			c.strips = append(c.strips, Strip{
				name: *light.Metadata.Name,
				id:   *light.Id,
			})
			c.logger.Info("found lightstrip", "name", *light.Metadata.Name, "id", *light.Id)
		default:
		}
	}
	return nil
}

func (c *HueClient) GetPlugs() []string {
	var plugs []string
	for name := range c.plugs {
		plugs = append(plugs, name)
	}
	sort.Strings(plugs)
	return plugs
}

func (c *HueClient) SetPlug(ctx context.Context, name string, state bool) error {
	plug, ok := c.plugs[name]
	if !ok {
		return fmt.Errorf("plug not found: %s", name)
	}
	action := LightPut{
		On: &On{
			On: &state,
		},
	}
	resp, err := c.Client.UpdateLight(ctx, plug.id, action)
	if err != nil {
		return fmt.Errorf("UpdateLight: %w", err)
	}
	_, err = ParseUpdateLightResponse(resp)
	if err != nil {
		return fmt.Errorf("ParseSetLightStateResponse: %w", err)
	}
	return nil
}

func pointerTo[T any](v T) *T {
	return &v
}

func gamut(c *color.XY) *GamutPosition {
	return &GamutPosition{
		X: pointerTo(c.X),
		Y: pointerTo(c.Y),
	}
}

func (c *HueClient) TurnOffStrip(ctx context.Context) error {
	for _, strip := range c.strips {
		action := LightPut{
			Type: pointerTo("light"),
			On: &On{
				On: pointerTo(false),
			},
		}
		resp, err := c.Client.UpdateLight(ctx, strip.id, action)
		if err != nil {
			return fmt.Errorf("UpdateLight: %w", err)
		}
		_, err = ParseUpdateLightResponse(resp)
		if err != nil {
			return fmt.Errorf("ParseSetLightStateResponse: %w", err)
		}
		c.logger.Info("turn off strip", "name", strip.name)
	}
	return nil
}

func (c *HueClient) SetStrip(ctx context.Context, cref string) error {
	if cref == "black" {
		return c.TurnOffStrip(ctx)
	}
	dimmedOn := pointerTo(float32(100.00000))
	// dimmedOff := pointerTo(float32(0.00001))
	xy, err := color.FindColorByName(cref)
	if err != nil {
		return fmt.Errorf("FindColorByName: %w", err)
	}
	for _, strip := range c.strips {
		action := LightPut{
			Type: pointerTo("light"),
			On: &On{
				On: pointerTo(true),
			},
			Color: &Color{
				Xy: gamut(xy),
			},
			Dimming: &Dimming{
				Brightness: dimmedOn,
			},
		}
		resp, err := c.Client.UpdateLight(ctx, strip.id, action)
		if err != nil {
			return fmt.Errorf("UpdateLight: %w", err)
		}
		_, err = ParseUpdateLightResponse(resp)
		if err != nil {
			return fmt.Errorf("ParseSetLightStateResponse: %w", err)
		}
		c.logger.Info("set strip", "name", strip.name, "cref", cref)
	}

	return nil
}
