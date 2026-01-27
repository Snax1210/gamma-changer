//go:build windows

package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"gamma-changer/internal/app"
)

// ShowHotkeyBindDialog 显示热键绑定对话框
func ShowHotkeyBindDialog(
	w fyne.Window,
	idx int,
	name string,
	status func(string),
	bindingHint *widget.Label,
	reloadCoreFromDisk func(),
	refreshPresetList func(),
) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Press A-Z or 0-9 (no need to hold Ctrl/Alt)")
	entry.Wrapping = fyne.TextWrapOff

	info := widget.NewLabel(fmt.Sprintf(
		"Binding preset '%s'\nResult will be: Ctrl+Alt+<Key>\n(Press one key: A-Z / 0-9)",
		name,
	))
	info.Wrapping = fyne.TextWrapWord

	var d dialog.Dialog
	completed := false

	entry.OnChanged = func(s string) {
		if s == "" || completed {
			return
		}

		// 验证输入
		if !validateHotkeyInput(s) {
			status("Bind: only A-Z / 0-9 supported.")
			entry.SetText("")
			return
		}

		r := []rune(s)[0]
		key := strings.ToUpper(string(r))
		spec := formatHotkeySpec(key)

		// 保存热键
		stDisk := app.LoadStateOrDefault()
		if stDisk.Hotkeys == nil {
			stDisk.Hotkeys = map[string]string{}
		}
		stDisk.Hotkeys[fmt.Sprintf("%d", idx)] = spec
		_ = app.SaveState(stDisk)

		// 重载并刷新UI
		reloadCoreFromDisk()
		completed = true
		status(fmt.Sprintf("Bound preset '%s' -> %s", name, spec))
		bindingHint.SetText("Hotkeys: click Bind then press a key (Ctrl+Alt+<Key>)")

		if d != nil {
			d.Hide()
		}
		refreshPresetList()
	}

	content := container.NewVBox(info, entry)

	d = dialog.NewCustomConfirm(
		"Bind Hotkey",
		"Cancel",
		"",
		content,
		func(_ bool) {
			if completed {
				return
			}
			status("Bind canceled.")
			bindingHint.SetText("Hotkeys: click Bind then press a key (Ctrl+Alt+<Key>)")
		},
		w,
	)

	bindingHint.SetText(fmt.Sprintf("Binding preset '%s'... (via dialog)", name))
	d.Show()

	// 延迟聚焦到输入框
	time.AfterFunc(10*time.Millisecond, func() {
		w.Canvas().Focus(entry)
	})
}

// validateHotkeyInput 验证热键输入
func validateHotkeyInput(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := []rune(s)[0]
	key := strings.ToUpper(string(r))
	return len(key) == 1 && ((key[0] >= '0' && key[0] <= '9') || (key[0] >= 'A' && key[0] <= 'Z'))
}

// formatHotkeySpec 格式化热键规范
func formatHotkeySpec(key string) string {
	return "Ctrl+Alt+" + key
}
