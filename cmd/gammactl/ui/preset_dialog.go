//go:build windows

package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"gamma-changer/internal/app"
)

// ShowNewPresetDialog 显示新建预设对话框
func ShowNewPresetDialog(
	w fyne.Window,
	core *app.App,
	status func(string),
	reloadCoreFromDisk func(),
	refreshPresetList func(),
) {
	entry := widget.NewEntry()
	entry.SetText("New Preset")

	d := dialog.NewForm("New Preset", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", entry),
		},
		func(ok bool) {
			if !ok {
				return
			}
			name := strings.TrimSpace(entry.Text)
			if name == "" {
				status("Preset name cannot be empty.")
				return
			}

			// 写磁盘：Presets append
			st2 := app.LoadStateOrDefault()
			st2.Presets = append(st2.Presets, app.Preset{
				Name:   name,
				Params: core.State().Params,
			})
			if st2.Hotkeys == nil {
				st2.Hotkeys = map[string]string{}
			}
			_ = app.SaveState(st2)

			// 立即生效：重载 + 刷新
			reloadCoreFromDisk()
			status("Preset created: " + name)
			refreshPresetList()
		}, w)
	d.Resize(fyne.NewSize(420, 160))
	d.Show()
}

// ShowRenamePresetDialog 显示重命名预设对话框
func ShowRenamePresetDialog(
	w fyne.Window,
	idx int,
	currentName string,
	status func(string),
	reloadCoreFromDisk func(),
	refreshPresetList func(),
) {
	entry := widget.NewEntry()
	entry.SetText(currentName)

	d := dialog.NewForm("Rename Preset", "OK", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", entry),
		},
		func(ok bool) {
			if !ok {
				return
			}
			newName := strings.TrimSpace(entry.Text)
			if newName == "" {
				status("Preset name cannot be empty.")
				return
			}

			st2 := app.LoadStateOrDefault()
			if idx < 0 || idx >= len(st2.Presets) {
				return
			}
			st2.Presets[idx].Name = newName
			_ = app.SaveState(st2)

			reloadCoreFromDisk()
			status("Preset renamed.")
			refreshPresetList()
		}, w)
	d.Resize(fyne.NewSize(420, 160))
	d.Show()
}

// ShowDeletePresetDialog 显示删除预设确认对话框
func ShowDeletePresetDialog(
	w fyne.Window,
	idx int,
	name string,
	status func(string),
	reloadCoreFromDisk func(),
	refreshPresetList func(),
	shiftHotkeysAfterDelete func(map[string]string, int) map[string]string,
) {
	d := dialog.NewConfirm("Delete Preset",
		"Delete preset '"+name+"' ?",
		func(ok bool) {
			if !ok {
				return
			}
			st2 := app.LoadStateOrDefault()
			if idx < 0 || idx >= len(st2.Presets) {
				return
			}
			st2.Presets = append(st2.Presets[:idx], st2.Presets[idx+1:]...)
			st2.Hotkeys = shiftHotkeysAfterDelete(st2.Hotkeys, idx)

			_ = app.SaveState(st2)

			reloadCoreFromDisk()
			status("Preset deleted: " + name)
			refreshPresetList()
		}, w)
	d.Show()
}
