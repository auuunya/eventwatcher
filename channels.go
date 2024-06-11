package eventwatcher

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	// https://learn.microsoft.com/zh-cn/windows/win32/api/winevt/
	modwevtapi             = syscall.MustLoadDLL("Wevtapi.dll")
	procEvtOpenChannelEnum = modwevtapi.MustFindProc("EvtOpenChannelEnum")
	procEvtNextChannelPath = modwevtapi.MustFindProc("EvtNextChannelPath")
	procEvtClose           = modwevtapi.MustFindProc("EvtClose")
)

func EvtOpenChannelEnum(session syscall.Handle) (syscall.Handle, error) {
	handle, _, err := procEvtOpenChannelEnum.Call(
		uintptr(session),
		0,
	)
	if handle == 0 {
		return InvalidHandle, errors.New("failed to open event channel enumeration handle: " + err.Error())
	}
	return syscall.Handle(handle), nil
}

func evtClose(handle syscall.Handle) error {
	return EvtClose(handle)
}

func EvtClose(handle syscall.Handle) error {
	ret, _, err := procEvtClose.Call(uintptr(handle))
	if ret == 0 {
		return errors.New("failed to close handle: " + err.Error())
	}
	return nil
}

func EvtNextChannelPath(handle syscall.Handle) ([]string, error) {
	channels := []string{}
	for {
		var bufferSize, bufferUsed uint32
		ret, _, err := procEvtNextChannelPath.Call(
			uintptr(handle),
			uintptr(bufferSize),
			0,
			uintptr(unsafe.Pointer(&bufferUsed)),
		)
		if ret == 0 {
			if err == ERROR_INSUFFICIENT_BUFFER {
				bufferSize = bufferUsed
				buffer := make([]uint16, bufferUsed)
				ret, _, _ := procEvtNextChannelPath.Call(
					uintptr(handle),
					uintptr(bufferSize),
					uintptr(unsafe.Pointer(&buffer[0])),
					uintptr(unsafe.Pointer(&bufferUsed)),
				)
				if ret == 0 {
					return nil, errors.New("failed to retrieve channel path after buffer allocation: " + err.Error())
				}
				channels = append(channels, syscall.UTF16ToString(buffer))
			} else if err == ERROR_NO_MORE_ITEMS {
				break
			} else {
				return nil, errors.New("failed to retrieve channel path: " + err.Error())
			}
		}
	}
	return channels, nil
}
