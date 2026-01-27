//go:build windows

package controller

import (
	"fmt"
	"time"

	"gamma-changer/internal/app"
)

// StateController 状态管理控制器
type StateController struct {
	core      *app.App
	hotkeyMgr *app.HotkeyManager
	status    func(string)
	debounce  time.Duration
	timer     *time.Timer
}

// NewStateController 创建状态管理控制器
func NewStateController(core *app.App, hotkeyMgr *app.HotkeyManager, statusFunc func(string)) *StateController {
	return &StateController{
		core:      core,
		hotkeyMgr: hotkeyMgr,
		status:    statusFunc,
		debounce:  80 * time.Millisecond,
	}
}

// ReloadFromDisk 从磁盘重载核心状态和热键
func (sc *StateController) ReloadFromDisk() {
	st := app.LoadStateOrDefault()
	sc.core = app.New(st)
	_ = sc.hotkeyMgr.ApplyFromState(sc.core.State(), sc.core)
}

// ApplyNow 立即应用当前参数
func (sc *StateController) ApplyNow() error {
	if err := sc.core.ApplyCurrent(); err != nil {
		sc.status("Apply failed: " + err.Error())
		return err
	}
	s := sc.core.State()
	sc.status(fmt.Sprintf("OK  γ=%.2f  b=%.2f  c=%.2f  (%s)",
		s.Params.Gamma, s.Params.Brightness, s.Params.Contrast, s.SelectedDisplay))
	_ = app.SaveState(s)
	return nil
}

// ScheduleApply 调度延迟应用（防抖）
func (sc *StateController) ScheduleApply() {
	if sc.timer != nil {
		sc.timer.Stop()
	}
	sc.timer = time.AfterFunc(sc.debounce, func() {
		err := sc.ApplyNow()
		if err != nil {
			return
		}
	})
}

// UpdateStatus 更新状态显示
func (sc *StateController) UpdateStatus(msg string) {
	if sc.status != nil {
		sc.status(msg)
	}
}

// GetCore 获取核心应用实例
func (sc *StateController) GetCore() *app.App {
	return sc.core
}

// UpdateParams 更新参数
func (sc *StateController) UpdateParams(updateFunc func(p *app.Params)) {
	sc.core.UpdateParams(updateFunc)
}
