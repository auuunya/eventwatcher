//go:build !windows
// +build !windows

package eventwatcher

// Non-Windows stubs for channel enumeration APIs.
func EvtOpenChannelEnum(session uintptr) (uintptr, error) { return 0, nil }
func EvtClose(handle uintptr) error                                   { return nil }
func EvtNextChannelPath(handle uintptr) ([]string, error)             { return nil, nil }
