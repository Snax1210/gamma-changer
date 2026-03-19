//go:build !singlekey

package app

// HotkeyPrefix returns the modifier prefix for hotkey display and binding.
// Default: Ctrl+Alt
const HotkeyPrefix = "Ctrl+Alt"

// FormatHotkeySpec returns the full hotkey spec string for a given key.
// Example: "A" → "Ctrl+Alt+A"
func FormatHotkeySpec(key string) string {
	return HotkeyPrefix + "+" + key
}

// HotkeyBindHint returns the UI hint text for hotkey binding dialog.
func HotkeyBindHint() string {
	return "Hotkeys: click Bind then press a key (Ctrl+Alt+<Key>)"
}

// HotkeyBindDialogInfo returns the instructional text shown inside the bind dialog.
func HotkeyBindDialogInfo(presetName string) string {
	return "Binding preset '" + presetName + "'\nResult will be: Ctrl+Alt+<Key>\n(Press one key: A-Z / 0-9)"
}

// RequiredModifiers returns the modifiers that will be combined with the key.
// In default mode: Ctrl+Alt is always applied.
func RequiredModifiers() bool {
	return true
}
