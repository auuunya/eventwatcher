package eventwatcher

import (
	"syscall"
)

const (
	InvalidHandle = syscall.Handle(0)

	ERROR_HANDLE_EOF          syscall.Errno = 38
	ERROR_INSUFFICIENT_BUFFER syscall.Errno = 122
	ERROR_NO_MORE_ITEMS       syscall.Errno = 259
	NO_ERROR                                = 0
)

const (
	EVENTLOG_SUCCESS          = 0x0000
	EVENTLOG_ERROR_TYPE       = 0x0001
	EVENTLOG_WARNING_TYPE     = 0x0002
	EVENTLOG_INFORMATION_TYPE = 0x0004
	EVENTLOG_AUDIT_SUCCESS    = 0x0008
	EVENTLOG_AUDIT_FAILURE    = 0x0010
)

const (
	// https://learn.microsoft.com/zh-cn/windows/win32/api/winbase/nf-winbase-readeventloga
	EVENTLOG_SEEK_READ       = 0x0002
	EVENTLOG_SEQUENTIAL_READ = 0x0001

	EVENTLOG_FORWARDS_READ  = 0x0004
	EVENTLOG_BACKWARDS_READ = 0x0008
)
