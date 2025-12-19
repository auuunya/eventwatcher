//go:build windows
// +build windows

package eventwatcher_test

import (
	"context"
	"sync"
	"testing"

	"github.com/auuunya/eventwatcher"
	"golang.org/x/sys/windows"
)

var (
	channel = "Application"

	eventId     uint32 = 123432
	wantEventId uint32 = 123432
	message     string
)

func TestEventWatcher(t *testing.T) {
	ctx := context.Background()
	notify := eventwatcher.NewEventNotifier(ctx)
	defer notify.Close()

	err := notify.AddWatcher(channel)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	// Register the event source
	handle, err := eventwatcher.RegisterEventSource(nil, windows.StringToUTF16Ptr(channel))
	if err != nil {
		return
	}
	defer eventwatcher.DeregisterEventSource(handle)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ch := range notify.EventLogChannel {
			val := eventwatcher.ParseEventLogData(ch.Buffer)
			if val.EventID != wantEventId {
				t.Errorf("unable to read application event log, event id: %d, want event id: %d\n", val.EventID, wantEventId)
			}
			break
		}
	}()

	// Write an empty event
	message = "This is an event log message!"
	err = eventwatcher.ReportEvent(handle, eventwatcher.EVENTLOG_INFORMATION_TYPE, 0, eventId, nil, []string{message}, nil)
	if err != nil {
		return
	}

	wg.Wait()
}

func TestEmptyEventWatcher(t *testing.T) {
	ctx := context.Background()
	notify := eventwatcher.NewEventNotifier(ctx)
	defer notify.Close()

	err := notify.AddWatcher(channel)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	// Register the event source
	handle, err := eventwatcher.RegisterEventSource(nil, windows.StringToUTF16Ptr(channel))
	if err != nil {
		return
	}
	defer eventwatcher.DeregisterEventSource(handle)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ch := range notify.EventLogChannel {
			val := eventwatcher.ParseEventLogData(ch.Buffer)
			if val.EventID != 0 {
				t.Errorf("unable to read application event log, event id: %d, want event id: %d\n", val.EventID, 0)
			}
			break
		}
	}()

	// Write an empty event
	message = ""
	err = eventwatcher.ReportEvent(handle, eventwatcher.EVENTLOG_INFORMATION_TYPE, 0, 0, nil, []string{message}, nil)
	if err != nil {
		return
	}
	wg.Wait()
}
