package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/auuunya/eventwatcher"
)

func main() {
	// Usage: pass a file path to watch via environment variable or default to ./example.log
	path := os.Getenv("EW_PATH")
	if path == "" {
		path = "example.log"
	}

	ctx := context.Background()
	n := eventwatcher.NewEventNotifier(ctx)
	defer n.Close()

	if err := n.AddWatcher(path); err != nil {
		fmt.Printf("failed to add watcher: %v\n", err)
		return
	}

	go func() {
		for ev := range n.EventLogChannel {
			fmt.Printf("event: name=%s len=%d content=%q\n", ev.Name, len(ev.Buffer), string(ev.Buffer))
		}
	}()

	// Just run until interrupted
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// small grace period to flush
	time.Sleep(100 * time.Millisecond)
	fmt.Println("exiting")
}