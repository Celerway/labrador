package main

import (
	"context"
	"fmt"
	"github.com/celerway/labrador/broker"
	"log/slog"
	"os"
	"os/signal"
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
}

func run(ctx context.Context, output *os.File, args []string, env []string) error {
	lh := slog.NewJSONHandler(output, &slog.HandlerOptions{})
	logger := slog.New(lh)

	br := broker.New(logger)

	err := br.Run(ctx)
	if err != nil {
		return fmt.Errorf("broker.Run: %w", err)
	}
	return nil
}
