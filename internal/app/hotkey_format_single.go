//go:build singlekey

package app

// HotkeyPrefix returns the modifier prefix for hotkey display and binding.
// Single-key mode: no modifier prefix.
const HotkeyPrefix = ""

// FormatHotkeySpec returns the full hotkey spec string for a given key.
// Example: "A" → "A"
func FormatHotkeySpec(key string) string {
	return key
}

// HotkeyBindHint returns the UI hint text for hotkey binding dialog.
func HotkeyBindHint() string {
	return "Hotkeys: click Bind then press a key (single-key mode)"
}

// HotkeyBindDialogInfo returns the instructional text shown inside the bind dialog.
func HotkeyBindDialogInfo(presetName string) string {
	return "Binding preset '" + presetName + "'\nSingle-key mode: just press one key\n(Press one key: A-Z / 0-9)"
}

// RequiredModifiers returns whether modifiers (Ctrl+Alt) should be combined with the key.
// In single-key mode: no modifiers required.
func RequiredModifiers() bool {
	return false
}
