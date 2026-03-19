//go:build windows

package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"gamma-changer/internal/app"
)

// PresetListUI 预设列表UI
type PresetListUI struct {
	Container    *fyne.Container
	BindingHint  *widget.Label
	status       func(string)
	reloadCore   func()
	refreshList  func()
	shiftHotkeys func(map[string]string, int) map[string]string
	sliders      *SlidersContainer
	window       fyne.Window
}

// NewPresetListUI 创建预设列表UI
func NewPresetListUI(
	status func(string),
	reloadCore func(),
	shiftHotkeys func(map[string]string, int) map[string]string,
	sliders *SlidersContainer,
	window fyne.Window,
) *PresetListUI {
	plui := &PresetListUI{
		Container:    container.NewVBox(),
		BindingHint:  widget.NewLabel(app.HotkeyBindHint()),
		status:       status,
		reloadCore:   reloadCore,
		shiftHotkeys: shiftHotkeys,
		sliders:      sliders,
		window:       window,
	}
	plui.BindingHint.Wrapping = fyne.TextWrapWord

	// 设置刷新函数
	plui.refreshList = plui.Refresh

	return plui
}

// Refresh 刷新预设列表
func (plui *PresetListUI) Refresh() {
	plui.Container.Objects = nil

	// 每次都从磁盘读，保证新增/改名/删除立刻体现在UI
	stDisk := app.LoadStateOrDefault()

	for i, p := range stDisk.Presets {
		idx := i
		name := p.Name

		lblName := widget.NewLabel(name)
		lblHK := widget.NewLabel(plui.getHotkeyLabel(idx))

		// Apply按钮
		btnPresetApply := widget.NewButton("Apply", func() {
			plui.reloadCore()

			core := app.New(app.LoadStateOrDefault())
			if err := core.ApplyPreset(idx); err != nil {
				plui.status("Preset failed: " + err.Error())
				return
			}

			// 同步滑块
			s2 := core.State()
			plui.sliders.SyncFromState(s2.Params)

			plui.status("Preset applied: " + name)
			_ = app.SaveState(s2)
		})

		// Bind按钮
		btnBind := widget.NewButton("Bind", func() {
			ShowHotkeyBindDialog(
				plui.window,
				idx,
				name,
				plui.status,
				plui.BindingHint,
				plui.reloadCore,
				plui.refreshList,
			)
		})

		// Clear按钮
		btnClear := widget.NewButton("Clear", func() {
			st2 := app.LoadStateOrDefault()
			if st2.Hotkeys == nil {
				st2.Hotkeys = map[string]string{}
			}
			delete(st2.Hotkeys, strconv.Itoa(idx))
			_ = app.SaveState(st2)

			plui.reloadCore()
			plui.status("Cleared hotkey for: " + name)
			plui.Refresh()
		})

		// Save按钮
		btnSave := widget.NewButton("Save", func() {
			core := app.New(app.LoadStateOrDefault())
			cur := core.State().Params

			st2 := app.LoadStateOrDefault()
			if idx < 0 || idx >= len(st2.Presets) {
				return
			}
			st2.Presets[idx].Params = cur
			_ = app.SaveState(st2)

			plui.reloadCore()
			plui.status("Preset saved from current: " + name)
			plui.Refresh()
		})

		// Rename按钮
		btnRename := widget.NewButton("Rename", func() {
			ShowRenamePresetDialog(
				plui.window,
				idx,
				name,
				plui.status,
				plui.reloadCore,
				plui.refreshList,
			)
		})

		// Delete按钮
		btnDelete := widget.NewButton("Delete", func() {
			ShowDeletePresetDialog(
				plui.window,
				idx,
				name,
				plui.status,
				plui.reloadCore,
				plui.refreshList,
				plui.shiftHotkeys,
			)
		})

		line := container.NewGridWithColumns(8,
			lblName,
			lblHK,
			btnPresetApply,
			btnBind,
			btnClear,
			btnSave,
			btnRename,
			btnDelete,
		)
		plui.Container.Add(line)
	}

	plui.Container.Refresh()
}

// getHotkeyLabel 获取热键标签
func (plui *PresetListUI) getHotkeyLabel(idx int) string {
	stDisk := app.LoadStateOrDefault()
	if stDisk.Hotkeys == nil {
		return "-"
	}
	if v, ok := stDisk.Hotkeys[strconv.Itoa(idx)]; ok && v != "" {
		return v
	}
	return "-"
}

// GetNewPresetButton 获取新建预设按钮
func (plui *PresetListUI) GetNewPresetButton(core *app.App) *widget.Button {
	return widget.NewButton("New Preset (from current)", func() {
		ShowNewPresetDialog(
			plui.window,
			core,
			plui.status,
			plui.reloadCore,
			plui.refreshList,
		)
	})
}
