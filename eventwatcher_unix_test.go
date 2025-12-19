//go:build !windows
// +build !windows

package eventwatcher

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func TestEventWatcherUnixFile(t *testing.T) {
	f, err := os.CreateTemp("", "ew_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	ctx := context.Background()
	n := NewEventNotifier(ctx)
	defer n.Close()

	if err := n.AddWatcher(f.Name()); err != nil {
		t.Fatalf("AddWatcher failed: %v", err)
	}

	// Give the watcher a moment to start and register the path
	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case ch := <-n.EventLogChannel:
			if string(ch.Buffer) != "hello world" {
				t.Errorf("unexpected content: %q", string(ch.Buffer))
			}
		case <-time.After(5 * time.Second):
			t.Errorf("timed out waiting for file write event")
		}
	}()

	if err := os.WriteFile(f.Name(), []byte("hello world"), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	wg.Wait()
}
