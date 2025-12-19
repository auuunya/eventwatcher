//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/auuunya/eventwatcher"
	"golang.org/x/sys/windows"
)

func main() {
	channel := os.Getenv("EW_CHANNEL")
	if channel == "" {
		channel = "Application"
	}

	ctx := context.Background()
	n := eventwatcher.NewEventNotifier(ctx)
	defer n.Close()

	if err := n.AddWatcher(channel); err != nil {
		fmt.Printf("AddWatcher failed: %v\n", err)
		return
	}

	go func() {
		for ev := range n.EventLogChannel {
			r := eventwatcher.ParseEventLogData(ev.Buffer)
			fmt.Printf("name=%s id=%d content=%s\n", ev.Name, r.EventID, eventwatcher.FormatContent(ev.Buffer))
		}
	}()

	// register + report a test event so example shows output when run locally
	h, err := eventwatcher.RegisterEventSource(nil, windows.StringToUTF16Ptr(channel))
	if err == nil {
		_ = eventwatcher.ReportEvent(h, eventwatcher.EVENTLOG_INFORMATION_TYPE, 0, 1111, nil, []string{"example"}, nil)
		_ = eventwatcher.DeregisterEventSource(h)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("exiting")
}