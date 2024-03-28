package main

import (
	"context"
	"fmt"
	"github.com/celerway/labrador/broker"
	"github.com/celerway/labrador/web"
	"golang.org/x/sync/errgroup"
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
	fmt.Println("clean exit")
}

func run(ctx context.Context, output *os.File, args []string, env []string) error {
	lh := slog.NewJSONHandler(output, &slog.HandlerOptions{})
	logger := slog.New(lh)

	br := broker.New(logger)
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		return br.Run(ctx)
	})
	ws := web.New(":8080", br, logger)
	errGroup.Go(func() error {
		return ws.Run(ctx)
	})
	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("errgroup reported failure: %w", err)
	}
	return nil
}
