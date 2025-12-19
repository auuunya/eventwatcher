//go:build windows
// +build windows

package eventwatcher

import (
	"context"
	"syscall"
)

// Windows-specific methods implemented in this file.
// The EventWatcher struct is defined in eventwatcher_common.go.

// NewEventWatcher creates a new EventWatcher instance.
func NewEventWatcher(ctx context.Context, name string, eventChan chan *EventEntry) *EventWatcher {
	ctx, cancel := context.WithCancel(ctx)
	return &EventWatcher{
		Name:      name,
		ctx:       ctx,
		cancel:    cancel,
		eventChan: eventChan,
		stopCh:    make(chan struct{}),
	}
}

// Init initializes the EventWatcher instance.
func (ew *EventWatcher) Init() error {
	handle, err := openEventLog(ew.Name)
	if err != nil {
		return err
	}
	ew.handle = handle

	offset, err := eventRecordNumber(ew.handle)
	if err != nil {
		return err
	}
	ew.offset = offset

	if ew.eventHandle, err = createEvent(nil, 0, 1, nil); err != nil {
		return err
	}

	if ew.cancelHandle, err = createEvent(nil, 1, 0, nil); err != nil {
		return err
	}
	return nil
}

// Close cancels the context and triggers the cancel event.
func (ew *EventWatcher) Close() {
	ew.cancel()
	setEvent(ew.cancelHandle)
	close(ew.stopCh)
}

// CloseHandles closes all handles associated with the EventWatcher.
func (ew *EventWatcher) CloseHandles() error {
	var err error
	if ew.handle != 0 {
		if e := closeEventLog(ew.handle); e != nil {
			err = e
		}
	}
	if ew.cancelHandle != 0 {
		if e := closeHandle(ew.cancelHandle); e != nil {
			err = e
		}
	}
	if ew.eventHandle != 0 {
		if e := closeHandle(ew.eventHandle); e != nil {
			err = e
		}
	}
	return err
}

// Listen monitors the event log and processes changes.
func (ew *EventWatcher) Listen() {
	defer ew.CloseHandles()

	if err := notifyChange(ew.handle, ew.eventHandle); err != nil {
		return
	}

	for {
		select {
		case <-ew.stopCh:
			return
		case <-ew.ctx.Done():
			return
		default:
			handles := []syscall.Handle{ew.eventHandle, ew.cancelHandle}
			event, err := waitForMultipleObjects(handles, false, syscall.INFINITE)
			if err != nil {
				return
			}
			switch event {
			case syscall.WAIT_OBJECT_0:
				rn, err := eventRecordNumber(ew.handle)
				if err != nil {
					return
				}
				if ew.offset == rn {
					continue
				}
				ew.offset = rn

				buf, err := readEventLog(ew.handle, EVENTLOG_SEEK_READ|EVENTLOG_FORWARDS_READ, ew.offset)
				if err != nil {
					return
				}
				ew.eventChan <- &EventEntry{
					Name:   ew.Name,
					Handle: ew.handle,
					Buffer: buf,
				}

				if err := resetEvent(ew.eventHandle); err != nil {
					return
				}
				if err := resetEvent(ew.cancelHandle); err != nil {
					return
				}
			case syscall.WAIT_OBJECT_0 + 1:
				return
			default:
				return
			}
		}
	}
}
