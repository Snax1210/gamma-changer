//go:build windows

package ui

import (
	"fyne.io/fyne/v2/widget"

	"gamma-changer/internal/app"
)

// NewDisplaySelector 创建显示器选择器
func NewDisplaySelector(displayNames []string, core *app.App, scheduleApply func()) *widget.Select {
	displaySelect := widget.NewSelect(displayNames, func(s string) {
		core.SetSelectedDisplay(s)
		scheduleApply()
	})
	displaySelect.PlaceHolder = "Select display"

	// 设置当前选中的显示器
	if core.State().SelectedDisplay != "" {
		displaySelect.SetSelected(core.State().SelectedDisplay)
	}

	return displaySelect
}
