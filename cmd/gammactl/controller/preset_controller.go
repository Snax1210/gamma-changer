//go:build windows

package controller

import (
	"strconv"

	"gamma-changer/internal/app"
)

// PresetController 预设操作控制器
type PresetController struct {
	status func(string)
}

// NewPresetController 创建预设操作控制器
func NewPresetController(statusFunc func(string)) *PresetController {
	return &PresetController{
		status: statusFunc,
	}
}

// ApplyPreset 应用预设
func (pc *PresetController) ApplyPreset(core *app.App, idx int) error {
	if err := core.ApplyPreset(idx); err != nil {
		pc.status("Preset failed: " + err.Error())
		return err
	}
	pc.status("Preset applied")
	_ = app.SaveState(core.State())
	return nil
}

// SavePreset 保存预设
func (pc *PresetController) SavePreset(core *app.App, idx int) error {
	st := app.LoadStateOrDefault()
	if idx < 0 || idx >= len(st.Presets) {
		return nil
	}
	st.Presets[idx].Params = core.State().Params
	_ = app.SaveState(st)
	pc.status("Preset saved from current")
	return nil
}

// RenamePreset 重命名预设
func (pc *PresetController) RenamePreset(idx int, newName string) error {
	st := app.LoadStateOrDefault()
	if idx < 0 || idx >= len(st.Presets) {
		return nil
	}
	st.Presets[idx].Name = newName
	_ = app.SaveState(st)
	pc.status("Preset renamed.")
	return nil
}

// DeletePreset 删除预设
func (pc *PresetController) DeletePreset(idx int, name string) error {
	st := app.LoadStateOrDefault()
	if idx < 0 || idx >= len(st.Presets) {
		return nil
	}
	st.Presets = append(st.Presets[:idx], st.Presets[idx+1:]...)
	st.Hotkeys = pc.ShiftHotkeysAfterDelete(st.Hotkeys, idx)
	_ = app.SaveState(st)
	pc.status("Preset deleted: " + name)
	return nil
}

// CreatePreset 创建新预设
func (pc *PresetController) CreatePreset(core *app.App, name string) error {
	st := app.LoadStateOrDefault()
	st.Presets = append(st.Presets, app.Preset{
		Name:   name,
		Params: core.State().Params,
	})
	if st.Hotkeys == nil {
		st.Hotkeys = map[string]string{}
	}
	_ = app.SaveState(st)
	pc.status("Preset created: " + name)
	return nil
}

// BindHotkey 绑定热键
func (pc *PresetController) BindHotkey(idx int, spec string) error {
	st := app.LoadStateOrDefault()
	if st.Hotkeys == nil {
		st.Hotkeys = map[string]string{}
	}
	st.Hotkeys[strconv.Itoa(idx)] = spec
	_ = app.SaveState(st)
	return nil
}

// ClearHotkey 清除热键
func (pc *PresetController) ClearHotkey(idx int, name string) error {
	st := app.LoadStateOrDefault()
	if st.Hotkeys == nil {
		st.Hotkeys = map[string]string{}
	}
	delete(st.Hotkeys, strconv.Itoa(idx))
	_ = app.SaveState(st)
	pc.status("Cleared hotkey for: " + name)
	return nil
}

// ShiftHotkeysAfterDelete 删除预设后调整热键索引
func (pc *PresetController) ShiftHotkeysAfterDelete(hk map[string]string, deletedIdx int) map[string]string {
	out := map[string]string{}
	if hk == nil {
		return out
	}
	for k, v := range hk {
		i, e := strconv.Atoi(k)
		if e != nil {
			continue
		}
		if i == deletedIdx {
			continue
		}
		if i > deletedIdx {
			out[strconv.Itoa(i-1)] = v
		} else {
			out[strconv.Itoa(i)] = v
		}
	}
	return out
}

// GetHotkeyLabel 获取热键标签
func (pc *PresetController) GetHotkeyLabel(idx int) string {
	st := app.LoadStateOrDefault()
	if st.Hotkeys == nil {
		return "-"
	}
	if v, ok := st.Hotkeys[strconv.Itoa(idx)]; ok && v != "" {
		return v
	}
	return "-"
}
