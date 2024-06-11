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

type EventLogType int

const (
	EVENTLOG_ERROR_TYPE       EventLogType = 0x0001 //错误事件
	EVENTLOG_AUDIT_FAILURE    EventLogType = 0x0010 //失败审核事件
	EVENTLOG_AUDIT_SUCCESS    EventLogType = 0x0008 //成功审核事件
	EVENTLOG_INFORMATION_TYPE EventLogType = 0x0004 //信息事件
	EVENTLOG_WARNING_TYPE     EventLogType = 0x0002 //警告事件
)

// https://learn.microsoft.com/zh-cn/windows/win32/api/winbase/nf-winbase-readeventloga
const (
	EVENTLOG_SEEK_READ       = 0x0002
	EVENTLOG_SEQUENTIAL_READ = 0x0001

	EVENTLOG_FORWARDS_READ  = 0x0004
	EVENTLOG_BACKWARDS_READ = 0x0008
)
