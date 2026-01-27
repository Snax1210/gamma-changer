//go:build windows

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// BuildMainLayout 构建主布局
func BuildMainLayout(
	w fyne.Window,
	displaySelect *widget.Select,
	presetListUI *PresetListUI,
	sliders *SlidersContainer,
	status *widget.Label,
	btnApply *widget.Button,
	btnReset *widget.Button,
	btnNewPreset *widget.Button,
) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Gamma / Brightness / Contrast (Software LUT, Real-time)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Display"),
		displaySelect,

		widget.NewSeparator(),
		widget.NewLabelWithStyle("Presets & Hotkeys (Global)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnNewPreset),
		presetListUI.BindingHint,
		presetListUI.Container,

		widget.NewSeparator(),
		CreateSliderRow("Gamma", sliders.GammaSlider, sliders.GammaValue),
		CreateSliderRow("Brightness", sliders.BrightSlider, sliders.BrightValue),
		CreateSliderRow("Contrast", sliders.ContrastSlider, sliders.ContrastValue),

		container.NewHBox(btnApply, btnReset),
		widget.NewSeparator(),
		status,
	)
}

// SetWindowContent 设置窗口内容
func SetWindowContent(w fyne.Window, content fyne.CanvasObject) {
	w.SetContent(content)
}
