//go:build !windows
// +build !windows

package eventwatcher

// Non-Windows minimal parser stubs.
// On non-Windows platforms we don't have EventLog structures; these
// functions provide safe fallbacks for compilation and will be
// implemented properly later when adding Unix event support.

type EventLogRecord struct{}

func ParseEventLogData(buf []byte) *EventLogRecord {
	return &EventLogRecord{}
}

func ParserEventLogData(buf []byte) (*EventLogRecord, error) {
	return &EventLogRecord{}, nil
}

func FormatContent(buf []byte) string {
	return ""
}

func FormatMessage(errorCode uint32) string {
	return ""
}

func LookupAccountSid(buf []byte, sidlen, sidoffset uint32) (string, string, error) {
	return "", "", nil
}