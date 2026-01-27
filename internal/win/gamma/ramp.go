//go:build windows

package gamma

import (
	"errors"
	"fmt"
	"math"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modgdi32               = windows.NewLazySystemDLL("gdi32.dll")
	procCreateDCA          = modgdi32.NewProc("CreateDCA")
	procDeleteDC           = modgdi32.NewProc("DeleteDC")
	procGetDeviceGammaRamp = modgdi32.NewProc("GetDeviceGammaRamp")
	procSetDeviceGammaRamp = modgdi32.NewProc("SetDeviceGammaRamp")
)

type Ramp struct {
	Data [3 * 256]uint16
}

func (r *Ramp) ptr() unsafe.Pointer { return unsafe.Pointer(&r.Data[0]) }

type Params struct {
	Gamma      float64 // 1.0 neutral
	Brightness float64 // -1..1
	Contrast   float64 // 1.0 neutral

	RGain float64 // 1.0 neutral
	GGain float64
	BGain float64
}

func DefaultParams() Params {
	return Params{
		Gamma:      1.0,
		Brightness: 0.0,
		Contrast:   1.0,
		RGain:      1.0,
		GGain:      1.0,
		BGain:      1.0,
	}
}

// Display is "\\.\DISPLAY1"
func CreateDisplayDC(displayName string) (windows.Handle, error) {
	driver, _ := windows.BytePtrFromString("DISPLAY")
	device, err := windows.BytePtrFromString(displayName)
	if err != nil {
		return 0, err
	}
	r1, _, e1 := procCreateDCA.Call(
		uintptr(unsafe.Pointer(driver)),
		uintptr(unsafe.Pointer(device)),
		0,
		0,
	)
	if r1 == 0 {
		if e1 != nil && e1 != syscall.Errno(0) {
			return 0, fmt.Errorf("CreateDCA failed: %v", e1)
		}
		return 0, errors.New("CreateDCA failed")
	}
	return windows.Handle(r1), nil
}

func DeleteDC(hdc windows.Handle) error {
	r1, _, e1 := procDeleteDC.Call(uintptr(hdc))
	if r1 == 0 {
		if e1 != nil && e1 != syscall.Errno(0) {
			return fmt.Errorf("DeleteDC failed: %v", e1)
		}
		return errors.New("DeleteDC failed")
	}
	return nil
}

func GetRamp(hdc windows.Handle) (Ramp, error) {
	var r Ramp
	ret, _, e1 := procGetDeviceGammaRamp.Call(uintptr(hdc), uintptr(r.ptr()))
	if ret == 0 {
		if e1 != nil && e1 != syscall.Errno(0) {
			return Ramp{}, fmt.Errorf("GetDeviceGammaRamp failed: %v", e1)
		}
		return Ramp{}, errors.New("GetDeviceGammaRamp failed")
	}
	return r, nil
}

func SetRamp(hdc windows.Handle, r Ramp) error {
	ret, _, e1 := procSetDeviceGammaRamp.Call(uintptr(hdc), uintptr(r.ptr()))
	if ret == 0 {
		if e1 != nil && e1 != syscall.Errno(0) {
			return fmt.Errorf("SetDeviceGammaRamp failed: %v", e1)
		}
		return errors.New("SetDeviceGammaRamp failed")
	}
	return nil
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
func clampU16(x float64) uint16 {
	if x < 0 {
		return 0
	}
	if x > 65535 {
		return 65535
	}
	return uint16(x + 0.5)
}

// y = ((x-0.5)*contrast + 0.5 + brightness)
// y = pow(clamp(y), 1/gamma)
// then apply per-channel gain.
func BuildRamp(p Params) (Ramp, error) {
	if p.Gamma <= 0 || math.IsNaN(p.Gamma) || math.IsInf(p.Gamma, 0) {
		return Ramp{}, errors.New("gamma must be > 0")
	}
	if p.Contrast <= 0 || math.IsNaN(p.Contrast) || math.IsInf(p.Contrast, 0) {
		return Ramp{}, errors.New("contrast must be > 0")
	}
	if p.RGain <= 0 || p.GGain <= 0 || p.BGain <= 0 {
		return Ramp{}, errors.New("RGB gain must be > 0")
	}

	invGamma := 1.0 / p.Gamma
	apply := func(x float64, gain float64) float64 {
		y := (x-0.5)*p.Contrast + 0.5 + p.Brightness
		y = clamp01(y)
		y = math.Pow(y, invGamma)
		y = clamp01(y * gain)
		return y
	}

	var r Ramp
	for i := 0; i < 256; i++ {
		x := float64(i) / 255.0
		rv := apply(x, p.RGain)
		gv := apply(x, p.GGain)
		bv := apply(x, p.BGain)

		r.Data[i] = clampU16(rv * 65535.0)
		r.Data[256+i] = clampU16(gv * 65535.0)
		r.Data[512+i] = clampU16(bv * 65535.0)
	}
	return r, nil
}
