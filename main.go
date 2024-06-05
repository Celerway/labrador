package main

import (
	"context"
	"fmt"
	"github.com/celerway/labrador/broker"
	"github.com/celerway/labrador/web"
	"github.com/hashicorp/mdns"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := run(ctx, os.Stdout, os.Args, os.Environ())
	if err != nil {
		fmt.Println("run error: ", err)
		os.Exit(1)
	}
	fmt.Println("clean exit")
}

func run(rootCtx context.Context, output *os.File, args []string, env []string) error {
	_ = godotenv.Load()
	lh := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: getLogLevel(),
	})
	logger := slog.New(lh)
	mqttPort := getEnvInt("MQTT_PORT", 1883)
	br := broker.New(fmt.Sprintf(":%d", mqttPort), logger)
	errGroup := new(errgroup.Group)
	ctx, cancel := context.WithCancel(rootCtx)
	errGroup.Go(func() error {
		defer cancel()
		return br.Run(ctx)
	})
	webAddr := getEnvString("WEB_ADDR", ":8080")
	downloadFolder := getEnvString("DOWNLOAD", "download")
	ws := web.New(webAddr, downloadFolder, br, logger)
	errGroup.Go(func() error {
		defer cancel()
		return ws.Run(ctx)
	})

	// announce the MQTT service
	server, err := announceName(mqttPort, logger)
	if err != nil {
		cancel()
		return fmt.Errorf("announceName: %w", err)
	}
	defer server.Shutdown()

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("errgroup reported failure: %w", err)
	}
	return nil
}

func getEnvInt(key string, defaultValue int) int {
	str, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	i, err := fmt.Sscanf(str, "%d", &defaultValue)
	if err != nil {
		return defaultValue
	}
	return i
}

func getEnvString(s string, s2 string) string {
	str, ok := os.LookupEnv(s)
	if !ok {
		return s2
	}
	return str
}

// getLogLevel returns the log level from the environment variable LOG_LEVEL.
func getLogLevel() slog.Level {
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func announceName(port int, logger *slog.Logger) (*mdns.Server, error) {
	serviceName := "labrador"
	serviceType := "_mqtt._tcp"
	serviceInfo := []string{"path=/mqtt"}

	// Get the local IP address to advertise the service
	ips, err := getLocalIPs()
	if err != nil {
		return nil, fmt.Errorf("getLocalIPs: %w", err)
	}

	// Create a new mDNS service
	service, err := mdns.NewMDNSService(serviceName, serviceType,
		"local.", "", port, ips, serviceInfo)
	if err != nil {
		return nil, fmt.Errorf("mdns.NewMDNSService: %w", err)
	}

	// Create an mDNS server using the service definition
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, fmt.Errorf("mdns.NewServer: %w", err)
	}

	logger.Info("MQTT Broker announced on mDNS.", "service name", serviceName, "service type", serviceType, "port", port, "local IPs", ips)
	// Keep the server running until the program is stopped
	return server, nil

}

// getLocalIPs attempts to determine a non-loopback address for the host
func getLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		// ignore docker interfaces
		if strings.HasPrefix(iface.Name, "docker") {
			continue
		}
		// ignore tailscale and other tunnel interfaces on macs
		if strings.HasPrefix(iface.Name, "utun") {
			continue
		}
		// same on linux:
		if strings.HasPrefix(iface.Name, "tailscale") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip)
		}
	}
	if len(ips) == 0 {
		return nil, os.ErrNotExist
	}
	return ips, nil
}
