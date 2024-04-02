package gohue

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// setup
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	ret := m.Run()
	// teardown
	os.Exit(ret)
}

func TestPlug(t *testing.T) {
	// test
	brUrl, ok := os.LookupEnv("HUE_BRIDGE")
	if !ok {
		t.Fatal("HUE_BRIDGE not set")
	}
	username, ok := os.LookupEnv("HUE_USERNAME")
	if !ok {
		t.Fatal("HUE_USERNAME not set")
	}
	fmt.Println(brUrl, username)
	c, err := New(brUrl, username, makeDebugLogger())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = c.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("loaded")
	err = c.SetPlug(ctx, "POWER3", true)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	err = c.SetPlug(ctx, "POWER3", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStrip(t *testing.T) {
	// test
	brUrl, ok := os.LookupEnv("HUE_BRIDGE")
	if !ok {
		t.Fatal("HUE_BRIDGE not set")
	}
	username, ok := os.LookupEnv("HUE_USERNAME")
	if !ok {
		t.Fatal("HUE_USERNAME not set")
	}
	fmt.Println(brUrl, username)
	c, err := New(brUrl, username, makeDebugLogger())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = c.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("loaded")

	colors := []string{"red", "green", "blue", "black"}
	for _, color := range colors {
		err = c.SetStrip(ctx, color)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}

func makeDebugLogger() *slog.Logger {
	lh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(lh)
}
