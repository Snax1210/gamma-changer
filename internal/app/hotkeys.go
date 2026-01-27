//go:build windows

package app

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.design/x/hotkey"
)

type HotkeyManager struct {
	mu  sync.Mutex
	hks map[int]*hotkey.Hotkey
}

func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{hks: map[int]*hotkey.Hotkey{}}
}

func (m *HotkeyManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, hk := range m.hks {
		_ = hk.Unregister()
	}
	m.hks = map[int]*hotkey.Hotkey{}
}

func parseHotkeySpec(spec string) ([]hotkey.Modifier, hotkey.Key, error) {
	// 例："Ctrl+Alt+1" / "Ctrl+Alt+N"
	parts := strings.Split(spec, "+")
	if len(parts) == 0 {
		return nil, 0, fmt.Errorf("empty hotkey")
	}

	var mods []hotkey.Modifier
	keyPart := parts[len(parts)-1]
	for _, p := range parts[:len(parts)-1] {
		switch strings.ToLower(strings.TrimSpace(p)) {
		case "ctrl", "control":
			mods = append(mods, hotkey.ModCtrl)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "win", "meta":
			mods = append(mods, hotkey.ModWin)
		default:
			return nil, 0, fmt.Errorf("unknown modifier: %s", p)
		}
	}

	kp := strings.ToUpper(strings.TrimSpace(keyPart))

	// 0-9
	if len(kp) == 1 && kp[0] >= '0' && kp[0] <= '9' {
		// hotkey.Key0...Key9
		n, _ := strconv.Atoi(kp)
		return mods, hotkey.Key(int(hotkey.Key0) + n), nil
	}

	// A-Z
	if len(kp) == 1 && kp[0] >= 'A' && kp[0] <= 'Z' {
		return mods, hotkey.Key(int(hotkey.KeyA) + int(kp[0]-'A')), nil
	}

	// 常用特殊键可以按需扩展：F1..F12、PageUp 等
	switch kp {
	case "PGUP":
		return mods, 0x21, nil
	case "PGDN":
		return mods, 0x22, nil
	default:
		return nil, 0, fmt.Errorf("unsupported key: %s", keyPart)
	}
}

const resetHotkey = "Ctrl+Alt+R"

func (m *HotkeyManager) ApplyFromState(st State, a *App) error {
	m.Clear()

	// st.Hotkeys: map[presetIndex]string
	for k, spec := range st.Hotkeys {
		idx, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		mods, key, err := parseHotkeySpec(spec)
		if err != nil {
			return fmt.Errorf("hotkey preset %d: %w", idx, err)
		}

		hk := hotkey.New(mods, key)
		if err := hk.Register(); err != nil {
			return fmt.Errorf("register %s for preset %d failed: %w", spec, idx, err)
		}

		m.mu.Lock()
		m.hks[idx] = hk
		m.mu.Unlock()

		// 监听触发
		go func(presetIdx int, h *hotkey.Hotkey) {
			for range h.Keydown() {
				_ = a.ApplyPreset(presetIdx)
				_ = SaveState(a.State())
			}
		}(idx, hk)
	}
	return nil
}
