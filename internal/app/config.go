package app

import (
	"encoding/json"
	"errors"
	"gamma-changer/internal/win/gamma"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func configDir() string {
	dir, _ := os.UserConfigDir()
	if dir == "" {
		dir = "."
	}
	p := filepath.Join(dir, "gammactl")
	_ = os.MkdirAll(p, 0700)
	return p
}

func configPath() string { return filepath.Join(configDir(), "config.json") }

func LoadStateOrDefault() State {
	// 默认给几个预设，足够日常一键切换
	def := State{
		Params:  DefaultParams(),
		Hotkeys: map[string]string{},
		Presets: []Preset{
			{Name: "Default", Params: DefaultParams()},
			{Name: "Office", Params: Params{Gamma: 1.05, Brightness: 0.02, Contrast: 1.05, RGain: 1, GGain: 1, BGain: 1}},
			{Name: "Night", Params: Params{Gamma: 1.20, Brightness: -0.02, Contrast: 0.95, RGain: 1, GGain: 0.95, BGain: 0.90}},
			{Name: "Coding", Params: Params{Gamma: 1.10, Brightness: 0.00, Contrast: 1.15, RGain: 1, GGain: 1, BGain: 1}},
		},
		AutoStart: false,
	}

	data, err := os.ReadFile(configPath())
	if err != nil {
		return def
	}
	var st State
	if json.Unmarshal(data, &st) != nil {
		return def
	}
	// 补默认
	if st.Params.Gamma == 0 {
		st.Params = def.Params
	}
	if len(st.Presets) == 0 {
		st.Presets = def.Presets
	}
	return st
}

func SaveState(st State) error {
	data, _ := json.MarshalIndent(st, "", "  ")
	return os.WriteFile(configPath(), data, 0600)
}

func backupPath(display string) string {
	safe := filepath.Base(display)
	return filepath.Join(configDir(), safe+"_backup_ramp.json")
}

func SaveBackupOnce(display string, hdc windows.Handle) error {
	if _, err := os.Stat(backupPath(display)); err == nil {
		return nil // 已存在
	}
	r, err := gamma.GetRamp(hdc)
	if err != nil {
		return err
	}
	data, _ := json.Marshal(r)
	return os.WriteFile(backupPath(display), data, 0600)
}

func RestoreBackup(display string, hdc windows.Handle) error {
	data, err := os.ReadFile(backupPath(display))
	if err != nil {
		// 没备份就恢复“中性 ramp”
		r, _ := gamma.BuildRamp(gamma.DefaultParams())
		return gamma.SetRamp(hdc, r)
	}
	var r gamma.Ramp
	if json.Unmarshal(data, &r) != nil {
		return errors.New("invalid backup ramp file")
	}
	return gamma.SetRamp(hdc, r)
}
