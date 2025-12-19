//go:build windows
// +build windows

package main

import (
	"context"
	"testing"
	"time"

	"github.com/auuunya/eventwatcher"
	"golang.org/x/sys/windows"
)

func TestWinWatchExample(t *testing.T) {
	channel := "Application"
	ctx := context.Background()
	n := eventwatcher.NewEventNotifier(ctx)
	defer n.Close()

	if err := n.AddWatcher(channel); err != nil {
		t.Skipf("AddWatcher failed (needs Windows Event Log): %v", err)
	}

	// register event source
	h, err := eventwatcher.RegisterEventSource(nil, windows.StringToUTF16Ptr(channel))
	if err != nil {
		t.Skipf("RegisterEventSource failed: %v", err)
	}
	defer eventwatcher.DeregisterEventSource(h)

	// report an event
	if err := eventwatcher.ReportEvent(h, eventwatcher.EVENTLOG_INFORMATION_TYPE, 0, 2222, nil, []string{"example from test"}, nil); err != nil {
		t.Fatalf("ReportEvent failed: %v", err)
	}

	// wait for the event
	select {
	case ev := <-n.EventLogChannel:
		r := eventwatcher.ParseEventLogData(ev.Buffer)
		if r.EventID != 2222 {
			t.Fatalf("unexpected event id: %d", r.EventID)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for windows event")
	}
}