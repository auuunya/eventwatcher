package eventwatcher

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type EventLogRecord struct {
	Length              uint32
	Reserved            uint32
	RecordNumber        uint32
	TimeGenerated       uint32
	TimeWritten         uint32
	EventID             uint32
	EventType           uint16
	NumStrings          uint16
	EventCategory       uint16
	ReservedFlags       uint16
	ClosingRecordNumber uint32
	StringOffset        uint32
	UserSidLength       uint32
	UserSidOffset       uint32
	DataLength          uint32
	DataOffset          uint32
}

// ParseEventLogData parses the event log data.
func ParseEventLogData(buf []byte) *EventLogRecord {
	var record EventLogRecord
	index := 0
	for {
		if index+int(unsafe.Sizeof(record)) > len(buf) {
			break
		}
		record = *(*EventLogRecord)(unsafe.Pointer(&buf[index]))
		index += int(record.Length)
		if index >= len(buf) {
			break
		}
	}
	return &record
}

func ParserEventLogData(buf []byte) (*EventLogRecord, error) {
	if len(buf) < int(unsafe.Sizeof(EventLogRecord{})) {
		return nil, windows.ERROR_INSUFFICIENT_BUFFER
	}
	record := (*EventLogRecord)(unsafe.Pointer(&buf[0]))
	return record, nil
}

func FormatContent(buf []byte) string {
	r := (*EventLogRecord)(unsafe.Pointer(&buf[0]))
	return syscall.UTF16ToString((*[1 << 10]uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + uintptr(r.StringOffset)))[:])
}

func FormatMessage(errorCode uint32) string {
	var messageBuffer [4096]uint16
	flags := uint32(windows.FORMAT_MESSAGE_FROM_SYSTEM | windows.FORMAT_MESSAGE_IGNORE_INSERTS)

	numChars, err := windows.FormatMessage(flags, 0, errorCode, 0, messageBuffer[:], nil)
	if err != nil {
		return fmt.Sprintf("Unknown error 0x%x", errorCode)
	}

	return windows.UTF16ToString(messageBuffer[:numChars])
}

// LookupAccountSid retrieves the account name and domain name for the specified SID.
func LookupAccountSid(buf []byte, sidlen, sidoffset uint32) (string, string, error) {
	var userSID *windows.SID
	if sidlen < 0 {
		return "", "", fmt.Errorf("unable to get sid")
	}
	userSID = (*windows.SID)(unsafe.Pointer(&buf[sidoffset]))
	var nameLen, domainLen uint32
	var sidType uint32

	// Initial call to determine the buffer sizes.
	err := windows.LookupAccountSid(nil, userSID, nil, &nameLen, nil, &domainLen, &sidType)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return "", "", err
	}

	// Allocate buffers.
	nameBuffer := make([]uint16, nameLen)
	domainBuffer := make([]uint16, domainLen)

	err = windows.LookupAccountSid(nil, userSID, &nameBuffer[0], &nameLen, &domainBuffer[0], &domainLen, &sidType)
	if err != nil {
		return "", "", err
	}

	return windows.UTF16ToString(nameBuffer), windows.UTF16ToString(domainBuffer), nil
}
