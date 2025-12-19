//go:build !windows
// +build !windows

package eventwatcher

import (
	"errors"
)

// Non-Windows stubs so package compiles on macOS/Linux. These should
// be replaced by real implementations (e.g., using fsnotify) in future
// work.

func openEventLog(name string) (uintptr, error) {
	return 0, errors.New("openEventLog not implemented on this OS")
}

func OpenEventLog(name string) (uintptr, error) {
	return 0, errors.New("OpenEventLog not implemented on this OS")
}

func closeEventLog(handle uintptr) error {
	return errors.New("closeEventLog not implemented on this OS")
}

func CloseEventLog(handle uintptr) error {
	return errors.New("CloseEventLog not implemented on this OS")
}

func notifyChange(handle, event uintptr) error {
	return errors.New("notifyChange not implemented on this OS")
}

func NotifyChangeEventLog(handle, event uintptr) error {
	return errors.New("NotifyChangeEventLog not implemented on this OS")
}

func eventRecordNumber(handle uintptr) (uint32, error) {
	return 0, errors.New("eventRecordNumber not implemented on this OS")
}

func EventLogRecordNumber(handle uintptr) (uint32, error) {
	return 0, errors.New("EventLogRecordNumber not implemented on this OS")
}

func readEventLog(handle uintptr, flags, offset uint32) ([]byte, error) {
	return nil, errors.New("readEventLog not implemented on this OS")
}

func ReadEventLog(handle uintptr, flags, offset uint32) ([]byte, error) {
	return nil, errors.New("ReadEventLog not implemented on this OS")
}

func createEvent(
	eventAttributes *uintptr,
	manualReset, initialState uint32,
	name *uint16,
) (uintptr, error) {
	return 0, errors.New("createEvent not implemented on this OS")
}

func CreateEvent(
	eventAttributes *uintptr,
	manualReset, initialState uint32,
	name *uint16,
) (uintptr, error) {
	return 0, errors.New("CreateEvent not implemented on this OS")
}

func resetEvent(handle uintptr) error {
	return errors.New("resetEvent not implemented on this OS")
}

func ResetEvent(handle uintptr) error {
	return errors.New("ResetEvent not implemented on this OS")
}

func waitForMultipleObjects(
	handles []uintptr,
	waitAll bool,
	waitMilliseconds uint32) (uint32, error) {
	return 0, errors.New("waitForMultipleObjects not implemented on this OS")
}

func WaitForMultipleObjects(
	handles []uintptr,
	waitAll bool,
	waitMilliseconds uint32,
) (event uint32, err error) {
	return 0, errors.New("WaitForMultipleObjects not implemented on this OS")
}

func setEvent(handle uintptr) error {
	return errors.New("setEvent not implemented on this OS")
}

func SetEvent(handle uintptr) error {
	return errors.New("SetEvent not implemented on this OS")
}

func closeHandle(handle uintptr) error {
	return errors.New("closeHandle not implemented on this OS")
}

func CloseHandle(handle uintptr) error {
	return errors.New("CloseHandle not implemented on this OS")
}

func registerEventSource(uncServerName, sourceName *uint16) (handle uintptr, err error) {
	return 0, errors.New("registerEventSource not implemented on this OS")
}

func RegisterEventSource(uncServerName, sourceName *uint16) (handle uintptr, err error) {
	return 0, errors.New("RegisterEventSource not implemented on this OS")
}

func reportEvent(log uintptr, etype uint16, category uint16, eventID uint32, userSid *uintptr, strings []string, binaryData []byte) error {
	return errors.New("reportEvent not implemented on this OS")
}

func ReportEvent(log uintptr, etype uint16, category uint16, eventID uint32, userSid *uintptr, strings []string, binaryData []byte) error {
	return errors.New("ReportEvent not implemented on this OS")
}

func deregisterEventSource(log uintptr) error {
	return errors.New("deregisterEventSource not implemented on this OS")
}

func DeregisterEventSource(log uintptr) error {
	return errors.New("DeregisterEventSource not implemented on this OS")
}