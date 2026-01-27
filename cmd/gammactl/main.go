//go:build windows

package main

import (
	"gamma-changer/cmd/gammactl/controller"
	"gamma-changer/cmd/gammactl/ui"
	"gamma-changer/internal/app"

	"fyne.io/fyne/v2/widget"
)

func main() {
	// 初始化核心应用
	core, _ := ui.InitializeCore()

	// 初始化热键管理器
	hotkeyMgr, _ := ui.InitializeHotkeyManager(core)

	// 初始化显示器列表
	displayNames, _ := ui.InitializeDisplays(core)

	// 创建Fyne应用和窗口
	gui := ui.NewApp()
	w := ui.NewWindow(gui)

	// 创建状态标签
	status := ui.NewStatusLabel()

	// 创建状态控制器
	stateCtrl := controller.NewStateController(core, hotkeyMgr, status.SetText)

	// 创建预设控制器
	presetCtrl := controller.NewPresetController(status.SetText)

	// 创建滑块容器
	sliders := ui.NewSlidersContainer(core, stateCtrl.UpdateParams, stateCtrl.ScheduleApply)

	// 创建显示器选择器
	displaySelect := ui.NewDisplaySelector(displayNames, core, stateCtrl.ScheduleApply)

	// 创建预设列表UI
	presetListUI := ui.NewPresetListUI(
		status.SetText,
		stateCtrl.ReloadFromDisk,
		presetCtrl.ShiftHotkeysAfterDelete,
		sliders,
		w,
	)

	// 创建新建预设按钮
	btnNewPreset := presetListUI.GetNewPresetButton(core)

	// 创建基础按钮
	btnReset := widget.NewButton("Reset (Restore backup)", func() {
		if err := core.Reset(); err != nil {
			status.SetText("Reset failed: " + err.Error())
			return
		}
		status.SetText("Reset OK")
		_ = app.SaveState(core.State())
	})

	btnApply := widget.NewButton("Apply Now", func() {
		err := stateCtrl.ApplyNow()
		if err != nil {
			return
		}
	})

	// 构建主布局
	content := ui.BuildMainLayout(
		w,
		displaySelect,
		presetListUI,
		sliders,
		status,
		btnApply,
		btnReset,
		btnNewPreset,
	)

	// 设置窗口内容
	ui.SetWindowContent(w, content)

	// 初次渲染预设列表
	presetListUI.Refresh()

	// 启动时应用一次
	stateCtrl.ScheduleApply()

	// 显示并运行窗口
	w.ShowAndRun()
}
