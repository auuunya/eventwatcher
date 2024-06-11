package eventwatcher

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syscall"
)

// EventNotifier manages a collection of EventWatchers.
type EventNotifier struct {
	EventLogChannel chan []byte
	watchers        map[string]*EventWatcher
	ctx             context.Context
	wg              sync.WaitGroup
	mu              sync.Mutex
}

// NewEventNotifier creates a new EventNotifier instance.
func NewEventNotifier(ctx context.Context) *EventNotifier {
	return &EventNotifier{
		ctx:             ctx,
		watchers:        make(map[string]*EventWatcher),
		EventLogChannel: make(chan []byte),
	}
}

// AddWatcher adds a new EventWatcher to the EventNotifier.
func (en *EventNotifier) AddWatcher(name string) error {
	en.mu.Lock()
	defer en.mu.Unlock()

	if _, exists := en.watchers[name]; exists {
		return errors.New(name + " event watcher already exists")
	}

	watcher := NewEventWatcher(en.ctx, name, en.EventLogChannel)
	if watcher == nil {
		return errors.New("unable to create event watcher")
	}

	if err := watcher.Init(); err != nil {
		return err
	}

	en.watchers[name] = watcher
	en.wg.Add(1)
	go func(watcher *EventWatcher) {
		defer en.wg.Done()
		watcher.Listen()
	}(watcher)
	return nil
}

// RemoveWatcher removes an EventWatcher from the EventNotifier.
func (en *EventNotifier) RemoveWatcher(name string) error {
	en.mu.Lock()
	defer en.mu.Unlock()

	watcher, exists := en.watchers[name]
	if !exists {
		return errors.New(name + " event watcher does not exist")
	}

	watcher.Close()
	delete(en.watchers, name)
	return nil
}

// Close shuts down all EventWatchers and waits for them to exit.
func (en *EventNotifier) Close() {
	en.mu.Lock()
	for _, watcher := range en.watchers {
		watcher.Close()
	}
	en.mu.Unlock()
	en.watchers = make(map[string]*EventWatcher)
	close(en.EventLogChannel)
	en.wg.Wait()
	fmt.Println("EventNotifier closed")
}

// GetWatcher retrieves an EventWatcher by name.
func (en *EventNotifier) GetWatcher(name string) (*EventWatcher, error) {
	en.mu.Lock()
	defer en.mu.Unlock()

	watcher, exists := en.watchers[name]
	if !exists {
		return nil, fmt.Errorf("%s event watcher not found", name)
	}
	return watcher, nil
}

// EventWatcher monitors an event log for changes.
type EventWatcher struct {
	Name         string
	handle       syscall.Handle
	offset       uint32
	eventHandle  syscall.Handle
	cancelHandle syscall.Handle
	ctx          context.Context
	cancel       context.CancelFunc
	eventChan    chan []byte
}

// NewEventWatcher creates a new EventWatcher instance.
func NewEventWatcher(ctx context.Context, name string, eventChan chan []byte) *EventWatcher {
	ctx, cancel := context.WithCancel(ctx)
	return &EventWatcher{
		Name:      name,
		ctx:       ctx,
		cancel:    cancel,
		eventChan: eventChan,
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

	eventHandle, err := createEvent(nil, 0, 1, nil)
	if err != nil {
		return err
	}
	ew.eventHandle = eventHandle

	cancelHandle, err := createEvent(nil, 1, 0, nil)
	if err != nil {
		return err
	}
	ew.cancelHandle = cancelHandle
	return nil
}

// Close cancels the context and triggers the cancel event.
func (ew *EventWatcher) Close() {
	// fmt.Printf("Watcher %s closing\n", ew.Name)
	ew.cancel()
	setEvent(ew.cancelHandle)
}

// CloseHandles closes all handles associated with the EventWatcher.
func (ew *EventWatcher) CloseHandles() error {
	var err error
	if ew.handle != 0 {
		err = closeEventLog(ew.handle)
	}
	if ew.cancelHandle != 0 {
		err = closeHandle(ew.cancelHandle)
	}
	if ew.eventHandle != 0 {
		err = closeHandle(ew.eventHandle)
	}
	return err
}

// Listen monitors the event log and processes changes.
func (ew *EventWatcher) Listen() {
	defer ew.CloseHandles()
	err := notifyChange(ew.handle, ew.eventHandle)
	if err != nil {
		return
	}
	for {
		fmt.Printf("333: %#v\n", ew.cancel)
		select {
		case <-ew.ctx.Done():
			// fmt.Printf("Watcher %s received cancel signal\n", ew.Name)
			// setEvent(ew.cancelHandle)
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
				ew.eventChan <- buf
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
