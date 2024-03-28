package main

import (
	"context"
	"fmt"
	"github.com/celerway/labrador/broker"
	"github.com/celerway/labrador/web"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log/slog"
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

func run(ctx context.Context, output *os.File, args []string, env []string) error {
	_ = godotenv.Load()
	lh := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: getLogLevel(),
	})
	logger := slog.New(lh)
	mqttAddr := getEnvString("MQTT_ADDR", ":1883")
	br := broker.New(mqttAddr, logger)
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		return br.Run(ctx)
	})
	webAddr := getEnvString("WEB_ADDR", ":8080")
	ws := web.New(webAddr, br, logger)
	errGroup.Go(func() error {
		return ws.Run(ctx)
	})
	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("errgroup reported failure: %w", err)
	}
	return nil
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
