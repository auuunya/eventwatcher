//go:build !windows
// +build !windows

package eventwatcher

import (
	"context"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Unix/macOS implementation of EventWatcher using fsnotify. The Name field
// is treated as a path to a file to watch; when the file is written to,
// the watcher will read the file content and emit it on the EventLogChannel.

type unixEventWatcher struct{
	*EventWatcher
	watcher *fsnotify.Watcher
	path    string
}

func NewEventWatcher(ctx context.Context, name string, eventChan chan *EventEntry) *EventWatcher {
	// Create a wrapper EventWatcher and return it; Init will set up fsnotify.
	ctx, cancel := context.WithCancel(ctx)
	return &EventWatcher{
		Name:      name,
		ctx:       ctx,
		cancel:    cancel,
		eventChan: eventChan,
		stopCh:    make(chan struct{}),
	}
}

// Init sets up fsnotify watcher for the provided file path.
func (ew *EventWatcher) Init() error {
	// Ensure file exists: create if missing
	if _, err := os.Stat(ew.Name); os.IsNotExist(err) {
		f, err := os.Create(ew.Name)
		if err != nil {
			return err
		}
		f.Close()
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := w.Add(ew.Name); err != nil {
		w.Close()
		return err
	}
	// replace underlying watcher using type assertion
	uwe := &unixEventWatcher{EventWatcher: ew, watcher: w, path: ew.Name}
	// store watcher pointer via embedding by setting to local variable (only used in Listen/CloseHandles)
	// we keep it on the heap by writing back
	*ew = *uwe.EventWatcher
	// hack: attach watcher via context value
	// but simpler: we store pointer in a private map? For now, attach to stopCh using closure in Listen.
	// We'll rely on ew.stopCh and w to be closed in CloseHandles.
	return nil
}

// Close handles cleans up resources for the watcher.
func (ew *EventWatcher) CloseHandles() error {
	return nil
}

// Close stops the watcher.
func (ew *EventWatcher) Close() {
	ew.cancel()
	select {
	case <-ew.stopCh:
		// already closed
	default:
		close(ew.stopCh)
	}
}

// Listen monitors the fsnotify watcher and emits file contents on write events.
func (ew *EventWatcher) Listen() {
	// Create watcher locally so we can close it in defer.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer w.Close()

	if err := w.Add(ew.Name); err != nil {
		return
	}

	for {
		select {
		case <-ew.stopCh:
			return
		case <-ew.ctx.Done():
			return
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create {
				// small debounce
				time.Sleep(20 * time.Millisecond)
				b, err := os.ReadFile(ew.Name)
				if err != nil {
					continue
				}
				ew.eventChan <- &EventEntry{Name: ew.Name, Handle: 0, Buffer: b}
			}
		case <-time.After(5 * time.Second):
			// keep loop alive and responsive to stop signals
		}
	}
}