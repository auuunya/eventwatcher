package eventwatcher

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type EventEntry struct {
	Name   string  `json:"name"`
	Handle uintptr `json:"handle"`
	Buffer []byte  `json:"buffer"`
}

// EventNotifier manages a collection of EventWatchers.
type EventNotifier struct {
	EventLogChannel chan *EventEntry
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
		EventLogChannel: make(chan *EventEntry),
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
