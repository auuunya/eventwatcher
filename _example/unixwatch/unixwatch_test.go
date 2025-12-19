//go:build !windows
// +build !windows

package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/auuunya/eventwatcher"
)

func TestUnixWatchExample(t *testing.T) {
	f, err := os.CreateTemp("", "example_unix_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	ctx := context.Background()
	n := eventwatcher.NewEventNotifier(ctx)
	defer n.Close()

	if err := n.AddWatcher(f.Name()); err != nil {
		t.Fatalf("AddWatcher failed: %v", err)
	}

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	if err := os.WriteFile(f.Name(), []byte("example payload"), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	select {
	case ev := <-n.EventLogChannel:
		if string(ev.Buffer) != "example payload" {
			t.Fatalf("unexpected payload: %q", string(ev.Buffer))
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}