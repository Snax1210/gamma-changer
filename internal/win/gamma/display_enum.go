//go:build windows

package gamma

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modUser32               = windows.NewLazySystemDLL("user32.dll")
	procEnumDisplayDevicesW = modUser32.NewProc("EnumDisplayDevicesW")
)

const (
	DISPLAY_DEVICE_ATTACHED_TO_DESKTOP = 0x00000001
)

// DISPLAY_DEVICEW
type displayDeviceW struct {
	Cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

func enumDisplayDevices(adapter *uint16, i uint32, dd *displayDeviceW, flags uint32) error {
	r1, _, _ := procEnumDisplayDevicesW.Call(
		uintptr(unsafe.Pointer(adapter)),
		uintptr(i),
		uintptr(unsafe.Pointer(dd)),
		uintptr(flags),
	)
	if r1 == 0 {
		// 关键：EnumDisplayDevicesW 枚举结束时，GetLastError 可能返回 0
		// 这种情况不是“参数错误”，应当当作“没有更多设备”
		errno := syscall.GetLastError()
		if errno == nil {
			return syscall.ERROR_NO_MORE_FILES
		}
		return errno
	}
	return nil
}

type Display struct {
	Name string // "\\.\DISPLAY1"
}

func ListDisplays() ([]Display, error) {
	var res []Display

	for i := uint32(0); ; i++ {
		var dd displayDeviceW
		dd.Cb = uint32(unsafe.Sizeof(dd))

		err := enumDisplayDevices(nil, i, &dd, 0)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return nil, fmt.Errorf("EnumDisplayDevicesW(%d): %w", i, err)
		}

		if dd.StateFlags&DISPLAY_DEVICE_ATTACHED_TO_DESKTOP == 0 {
			continue
		}

		name := windows.UTF16ToString(dd.DeviceName[:])
		if name == "" {
			continue
		}
		res = append(res, Display{Name: name})
	}

	if len(res) == 0 {
		return nil, errors.New("no attached displays found")
	}
	return res, nil
}
