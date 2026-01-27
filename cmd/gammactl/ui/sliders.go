//go:build windows

package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"gamma-changer/internal/app"
)

// SliderConfig 滑块配置
type SliderConfig struct {
	Title     string
	Min       float64
	Max       float64
	Step      float64
	Initial   float64
	OnChanged func(float64)
}

// NewSlider 创建通用滑块
func NewSlider(config SliderConfig) (*widget.Slider, *widget.Label) {
	slider := widget.NewSlider(config.Min, config.Max)
	slider.Step = config.Step
	slider.Value = config.Initial

	value := widget.NewLabel(fmt.Sprintf("%.2f", slider.Value))

	slider.OnChanged = func(v float64) {
		value.SetText(fmt.Sprintf("%.2f", v))
		if config.OnChanged != nil {
			config.OnChanged(v)
		}
	}

	return slider, value
}

// NewGammaSlider 创建Gamma滑块
func NewGammaSlider(initial float64, onChanged func(float64)) (*widget.Slider, *widget.Label) {
	return NewSlider(SliderConfig{
		Title:     "Gamma",
		Min:       0.30,
		Max:       4.40,
		Step:      0.01,
		Initial:   initial,
		OnChanged: onChanged,
	})
}

// NewBrightnessSlider 创建Brightness滑块
func NewBrightnessSlider(initial float64, onChanged func(float64)) (*widget.Slider, *widget.Label) {
	return NewSlider(SliderConfig{
		Title:     "Brightness",
		Min:       -1.00,
		Max:       1.00,
		Step:      0.01,
		Initial:   initial,
		OnChanged: onChanged,
	})
}

// NewContrastSlider 创建Contrast滑块
func NewContrastSlider(initial float64, onChanged func(float64)) (*widget.Slider, *widget.Label) {
	return NewSlider(SliderConfig{
		Title:     "Contrast",
		Min:       0.10,
		Max:       3.00,
		Step:      0.01,
		Initial:   initial,
		OnChanged: onChanged,
	})
}

// CreateSliderRow 创建滑块行布局
func CreateSliderRow(title string, slider *widget.Slider, value *widget.Label) fyne.CanvasObject {
	return container.NewBorder(nil, nil, widget.NewLabel(title), value, slider)
}

// SlidersContainer 滑块容器
type SlidersContainer struct {
	GammaSlider    *widget.Slider
	GammaValue     *widget.Label
	BrightSlider   *widget.Slider
	BrightValue    *widget.Label
	ContrastSlider *widget.Slider
	ContrastValue  *widget.Label
}

// NewSlidersContainer 创建滑块容器
func NewSlidersContainer(core *app.App, updateParams func(func(p *app.Params)), scheduleApply func()) *SlidersContainer {
	sc := &SlidersContainer{}

	// Gamma滑块
	sc.GammaSlider, sc.GammaValue = NewGammaSlider(
		core.State().Params.Gamma,
		func(v float64) {
			updateParams(func(p *app.Params) { p.Gamma = v })
			scheduleApply()
		},
	)

	// Brightness滑块
	sc.BrightSlider, sc.BrightValue = NewBrightnessSlider(
		core.State().Params.Brightness,
		func(v float64) {
			updateParams(func(p *app.Params) { p.Brightness = v })
			scheduleApply()
		},
	)

	// Contrast滑块
	sc.ContrastSlider, sc.ContrastValue = NewContrastSlider(
		core.State().Params.Contrast,
		func(v float64) {
			updateParams(func(p *app.Params) { p.Contrast = v })
			scheduleApply()
		},
	)

	return sc
}

// SyncFromState 从状态同步滑块值
func (sc *SlidersContainer) SyncFromState(params app.Params) {
	sc.GammaSlider.SetValue(params.Gamma)
	sc.BrightSlider.SetValue(params.Brightness)
	sc.ContrastSlider.SetValue(params.Contrast)
}
