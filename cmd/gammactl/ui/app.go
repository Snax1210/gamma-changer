//go:build windows

package ui

import (
	"gamma-changer/internal/app"
	"gamma-changer/internal/win/gamma"

	"fyne.io/fyne/v2"
	fapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// NewApp 创建Fyne应用实例
func NewApp() fyne.App {
	return fapp.New()
}

// NewWindow 创建主窗口
func NewWindow(app fyne.App) fyne.Window {
	w := app.NewWindow("Gamma Changer")
	w.Resize(fyne.NewSize(780, 640))
	return w
}

// InitializeCore 初始化核心应用实例
func InitializeCore() (*app.App, error) {
	st := app.LoadStateOrDefault()
	core := app.New(st)
	return core, nil
}

// InitializeHotkeyManager 初始化热键管理器
func InitializeHotkeyManager(core *app.App) (*app.HotkeyManager, error) {
	hkMgr := app.NewHotkeyManager()
	err := hkMgr.ApplyFromState(core.State(), core)
	return hkMgr, err
}

// InitializeDisplays 初始化显示器列表
func InitializeDisplays(core *app.App) ([]string, error) {
	ds, err := gamma.ListDisplays()
	if err != nil {
		return nil, err
	}

	// 设置默认显示器
	if core.State().SelectedDisplay == "" && len(ds) > 0 {
		core.SetSelectedDisplay(ds[0].Name)
	}

	// 提取显示器名称
	displayNames := make([]string, 0, len(ds))
	for _, d := range ds {
		displayNames = append(displayNames, d.Name)
	}

	return displayNames, nil
}

// NewStatusLabel 创建状态标签
func NewStatusLabel() *widget.Label {
	return widget.NewLabel("Ready")
}
