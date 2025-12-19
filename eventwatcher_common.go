package eventwatcher

import "context"

// Common EventWatcher fields shared across platforms.
// Platform-specific files implement the platform behavior (Init, Listen, CloseHandles).

type EventWatcher struct {
	Name         string
	handle       uintptr
	offset       uint32
	eventHandle  uintptr
	cancelHandle uintptr
	ctx          context.Context
	cancel       context.CancelFunc
	eventChan    chan *EventEntry
	stopCh       chan struct{}
}