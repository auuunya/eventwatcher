package eventwatcher

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// https://learn.microsoft.com/zh-cn/windows/win32/api/winbase/
	modadvapi32                    = syscall.MustLoadDLL("advapi32.dll")
	procOpenEventLog               = modadvapi32.MustFindProc("OpenEventLogW")
	procReadEventLog               = modadvapi32.MustFindProc("ReadEventLogW")
	procCloseEventLog              = modadvapi32.MustFindProc("CloseEventLog")
	procNotifyChangeEventLog       = modadvapi32.MustFindProc("NotifyChangeEventLog")
	procGetNumberOfEventLogRecords = modadvapi32.MustFindProc("GetNumberOfEventLogRecords")
	procRegisterEventSourceW       = modadvapi32.MustFindProc("RegisterEventSourceW")
	procReportEventW               = modadvapi32.MustFindProc("ReportEventW")
	procDeregisterEventSource      = modadvapi32.MustFindProc("DeregisterEventSource")
	// https://learn.microsoft.com/zh-cn/windows/win32/api/synchapi/
	modkernel32                = syscall.MustLoadDLL("Kernel32.dll")
	procCreateEvent            = modkernel32.MustFindProc("CreateEventW")
	procResetEvent             = modkernel32.MustFindProc("ResetEvent")
	procWaitForMultipleObjects = modkernel32.MustFindProc("WaitForMultipleObjects")
	procSetEvent               = modkernel32.MustFindProc("SetEvent")
	procCloseHandle            = modkernel32.MustFindProc("CloseHandle")
)

func openEventLog(name string) (syscall.Handle, error) {
	return OpenEventLog(name)
}

func OpenEventLog(name string) (syscall.Handle, error) {
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return InvalidHandle, err
	}
	handle, _, err := procOpenEventLog.Call(
		0,
		uintptr(unsafe.Pointer(namePtr)),
	)
	if handle == 0 {
		return InvalidHandle, errors.New("failed to open event: " + err.Error())
	}
	return syscall.Handle(handle), nil
}

func closeEventLog(handle syscall.Handle) error {
	return CloseEventLog(handle)
}

func CloseEventLog(handle syscall.Handle) error {
	ret, _, err := procCloseEventLog.Call(uintptr(handle))
	if ret == 0 {
		return errors.New("failed to close event: " + err.Error())
	}
	return nil
}
func notifyChange(handle, event syscall.Handle) error {
	return NotifyChangeEventLog(handle, event)
}

func NotifyChangeEventLog(handle, event syscall.Handle) error {
	ret, _, err := procNotifyChangeEventLog.Call(
		uintptr(handle),
		uintptr(event),
	)
	if ret == NO_ERROR {
		return errors.New("failed to notify change event: " + err.Error())
	}
	return nil
}

func eventRecordNumber(handle syscall.Handle) (uint32, error) {
	return EventLogRecordNumber(handle)
}

func EventLogRecordNumber(handle syscall.Handle) (uint32, error) {
	var retVal uint32
	ret, _, err := procGetNumberOfEventLogRecords.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return retVal, nil
	}
	return 0, errors.New("failed to get number of handle: " + err.Error())
}

func readEventLog(handle syscall.Handle, flags, offset uint32) ([]byte, error) {
	return ReadEventLog(handle, flags, offset)
}

func ReadEventLog(handle syscall.Handle, flags, offset uint32) ([]byte, error) {
	var BUFFER_SIZE = 4096
	buffer := make([]byte, BUFFER_SIZE)
	var bytesRead, minByteNeeded uint32
	for {
		ret, _, err := procReadEventLog.Call(
			uintptr(handle),
			uintptr(flags),
			uintptr(offset),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(BUFFER_SIZE),
			uintptr(unsafe.Pointer(&bytesRead)),     // 传递 bytesRead 的地址，以获取实际读取的字节数
			uintptr(unsafe.Pointer(&minByteNeeded)), // 不需要返回的记录数
		)
		if ret == 0 {
			if err == ERROR_HANDLE_EOF {
				break
			} else if err == ERROR_INSUFFICIENT_BUFFER {
				buffer = make([]byte, minByteNeeded)
				BUFFER_SIZE = int(minByteNeeded)
				continue
			} else {
				return nil, err
			}
		}
		return buffer[:bytesRead], nil
	}
	return nil, errors.New("unable to read event buffer")
}

func createEvent(
	eventAttributes *syscall.SecurityAttributes,
	manualReset, initialState uint32,
	name *uint16,
) (syscall.Handle, error) {
	return CreateEvent(eventAttributes, manualReset, initialState, name)
}

func CreateEvent(
	eventAttributes *syscall.SecurityAttributes,
	manualReset, initialState uint32,
	name *uint16,
) (syscall.Handle, error) {
	ret, _, err := procCreateEvent.Call(
		uintptr(unsafe.Pointer(eventAttributes)),
		uintptr(manualReset),
		uintptr(initialState),
		uintptr(unsafe.Pointer(name)),
	)
	if ret == 0 {
		return InvalidHandle, errors.New("unable to create event: " + err.Error())
	}
	return syscall.Handle(ret), nil
}
func resetEvent(handle syscall.Handle) error {
	return ResetEvent(handle)
}
func ResetEvent(handle syscall.Handle) error {
	ret, _, err := procResetEvent.Call(
		uintptr(handle),
	)
	if ret == 0 {
		return errors.New("unable to reset event: " + err.Error())
	}
	return nil
}

func waitForMultipleObjects(
	handles []syscall.Handle,
	waitAll bool,
	waitMilliseconds uint32) (uint32, error) {
	return WaitForMultipleObjects(handles, waitAll, waitMilliseconds)
}

func WaitForMultipleObjects(
	handles []syscall.Handle,
	waitAll bool,
	waitMilliseconds uint32,
) (event uint32, err error) {
	var ptr *syscall.Handle
	if len(handles) > 0 {
		ptr = &handles[0]
	}
	ret, _, e1 := procWaitForMultipleObjects.Call(
		uintptr(len(handles)),
		uintptr(unsafe.Pointer(ptr)),
		uintptr(*(*int32)(unsafe.Pointer(&waitAll))),
		uintptr(waitMilliseconds),
	)
	event = uint32(ret)
	if event == syscall.WAIT_FAILED {
		err = errors.New("unable to wait for multiple objects: " + e1.Error())
		return
	}
	return
}
func setEvent(handle syscall.Handle) error {
	return SetEvent(handle)
}
func SetEvent(handle syscall.Handle) error {
	ret, _, err := procSetEvent.Call(
		uintptr(handle),
	)
	if ret == 0 {
		return errors.New("unable to set event: " + err.Error())
	}
	return nil
}
func closeHandle(handle syscall.Handle) error {
	return CloseHandle(handle)
}

func CloseHandle(handle syscall.Handle) error {
	ret, _, err := procCloseHandle.Call(
		uintptr(handle),
	)
	if ret == 0 {
		return errors.New("unable to close handle: " + err.Error())
	}
	return nil
}

func registerEventSource(uncServerName, sourceName *uint16) (handle syscall.Handle, err error) {
	return RegisterEventSource(uncServerName, sourceName)
}

func RegisterEventSource(uncServerName, sourceName *uint16) (handle syscall.Handle, err error) {
	r0, _, e1 := procRegisterEventSourceW.Call(
		uintptr(unsafe.Pointer(uncServerName)),
		uintptr(unsafe.Pointer(sourceName)),
	)
	handle = syscall.Handle(r0)
	if handle == 0 {
		err = e1
	}
	return
}

func reportEvent(log syscall.Handle, etype uint16, category uint16, eventID uint32, userSid *windows.SID, strings []string, binaryData []byte) error {
	return ReportEvent(log, etype, category, eventID, userSid, strings, binaryData)
}

func ReportEvent(log syscall.Handle, etype uint16, category uint16, eventID uint32, userSid *windows.SID, strings []string, binaryData []byte) error {
	var pstrings **uint16
	if len(strings) > 0 {
		var strs []uintptr
		for _, s := range strings {
			strs = append(strs, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(s))))
		}
		pstrings = &(*[1 << 30]*uint16)(unsafe.Pointer(&strs[0]))[0]
	}
	var binaryDataPtr *byte
	if len(binaryData) > 0 {
		binaryDataPtr = &binaryData[0]
	}
	r1, _, e1 := procReportEventW.Call(
		uintptr(log),
		uintptr(etype),
		uintptr(category),
		uintptr(eventID),
		uintptr(unsafe.Pointer(userSid)),
		uintptr(len(strings)),
		uintptr(len(binaryData)),
		uintptr(unsafe.Pointer(pstrings)),
		uintptr(unsafe.Pointer(binaryDataPtr)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}
func deregisterEventSource(log syscall.Handle) error {
	return DeregisterEventSource(log)
}

func DeregisterEventSource(log syscall.Handle) error {
	r1, _, e1 := procDeregisterEventSource.Call(
		uintptr(log),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}
