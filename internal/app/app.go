package app

import (
	"fmt"
	"gamma-changer/internal/win/gamma"
	"strconv"
	"sync"
)

type Params struct {
	Gamma      float64 `json:"gamma"`
	Brightness float64 `json:"brightness"`
	Contrast   float64 `json:"contrast"`
	RGain      float64 `json:"r_gain"`
	GGain      float64 `json:"g_gain"`
	BGain      float64 `json:"b_gain"`
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

type Preset struct {
	Name   string `json:"name"`
	Params Params `json:"params"`
}

type State struct {
	SelectedDisplay string   `json:"selected_display"`
	Params          Params   `json:"params"`
	Presets         []Preset `json:"presets"`
	AutoStart       bool     `json:"auto_start"`
	// 新增：preset 索引 -> 热键，例如 {"0":"Ctrl+Alt+1","2":"Ctrl+Alt+N"}
	Hotkeys map[string]string `json:"hotkeys"`
}

type App struct {
	mu    sync.Mutex
	state State
}

func New(initial State) *App { return &App{state: initial} }

func (a *App) State() State {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

func (a *App) SetSelectedDisplay(name string) {
	a.mu.Lock()
	a.state.SelectedDisplay = name
	a.mu.Unlock()
}

func (a *App) ApplyCurrent() error {
	a.mu.Lock()
	display := a.state.SelectedDisplay
	p := a.state.Params
	a.mu.Unlock()

	if display == "" {
		return fmt.Errorf("no display selected")
	}

	hdc, err := gamma.CreateDisplayDC(display)
	if err != nil {
		return err
	}
	defer gamma.DeleteDC(hdc)

	// 备份：只在首次没有备份时保存（你之前 CLI 那套备份逻辑也可以搬过来）
	// 这里简单化：每次启动时先读一次并保存到文件（见 config.go 里的 SaveBackupOnce）
	if err := SaveBackupOnce(display, hdc); err != nil {
		// 不致命：有些机器读不出来，但写可能还能写；你可以按需改为致命
	}

	ramp, err := gamma.BuildRamp(gamma.Params{
		Gamma:      p.Gamma,
		Brightness: p.Brightness,
		Contrast:   p.Contrast,
		RGain:      p.RGain,
		GGain:      p.GGain,
		BGain:      p.BGain,
	})
	if err != nil {
		return err
	}
	return gamma.SetRamp(hdc, ramp)
}

func (a *App) Reset() error {
	a.mu.Lock()
	display := a.state.SelectedDisplay
	a.mu.Unlock()
	if display == "" {
		return fmt.Errorf("no display selected")
	}
	hdc, err := gamma.CreateDisplayDC(display)
	if err != nil {
		return err
	}
	defer gamma.DeleteDC(hdc)
	return RestoreBackup(display, hdc)
}

func (a *App) AdjustBrightness(delta float64) {
	a.mu.Lock()
	a.state.Params.Brightness = clamp(a.state.Params.Brightness+delta, -1.0, 1.0)
	a.mu.Unlock()
}
func (a *App) AdjustContrast(delta float64) {
	a.mu.Lock()
	a.state.Params.Contrast = clamp(a.state.Params.Contrast+delta, 0.10, 3.00)
	a.mu.Unlock()
}
func (a *App) AdjustGamma(delta float64) {
	a.mu.Lock()
	a.state.Params.Gamma = clamp(a.state.Params.Gamma+delta, 0.30, 4.40)
	a.mu.Unlock()
}

func (a *App) ApplyPreset(i int) error {
	a.mu.Lock()
	if i < 0 || i >= len(a.state.Presets) {
		a.mu.Unlock()
		return fmt.Errorf("preset index out of range")
	}
	a.state.Params = a.state.Presets[i].Params
	a.mu.Unlock()
	return a.ApplyCurrent()
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// UpdateParams updates the current Params under lock.
// It returns a snapshot State after the update (useful for UI sync).
func (a *App) UpdateParams(fn func(p *Params)) State {
	a.mu.Lock()
	fn(&a.state.Params)
	st := a.state
	a.mu.Unlock()
	return st
}

func (a *App) AddPreset(name string) State {
	a.mu.Lock()
	defer a.mu.Unlock()

	p := Preset{
		Name:   name,
		Params: a.state.Params,
	}
	a.state.Presets = append(a.state.Presets, p)
	return a.state
}

func (a *App) UpdatePresetFromCurrent(i int) State {
	a.mu.Lock()
	defer a.mu.Unlock()

	if i < 0 || i >= len(a.state.Presets) {
		return a.state
	}
	a.state.Presets[i].Params = a.state.Params
	return a.state
}

func (a *App) DeletePreset(i int) State {
	a.mu.Lock()
	defer a.mu.Unlock()

	if i < 0 || i >= len(a.state.Presets) {
		return a.state
	}
	a.state.Presets = append(a.state.Presets[:i], a.state.Presets[i+1:]...)
	delete(a.state.Hotkeys, strconv.Itoa(i)) // 顺手移除绑定
	return a.state
}
