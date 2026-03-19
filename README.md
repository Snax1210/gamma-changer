# Gamma Changer

A powerful display calibration tool designed specifically for Windows platform that allows real-time adjustment of
monitor Gamma, brightness, contrast, and RGB gain through software LUT (Look-Up Table).

![Go Version](https://img.shields.io/badge/Go-1.24.10-blue)
![Platform](https://img.shields.io/badge/Platform-Windows-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## Preview

<img width="1321" height="931" alt="image" src="https://github.com/user-attachments/assets/3d329cef-3e57-453e-96ca-4dc781207d0a" />

## Features

- **Real-time Display Adjustment**: Adjust Gamma, brightness, contrast, and RGB gain values in real-time
- **Preset Management**: Create, save, apply, rename, and delete custom display presets
- **Flexible Hotkey Modes**: Choose between `Ctrl+Alt+Key` (default) or single-key hotkey mode at build time for quick
  preset switching
- **Multi-Monitor Support**: Manage settings for multiple displays independently
- **Backup & Restore**: Automatically backup original settings and restore them when needed
- **Auto-Start**: Option to launch the application automatically on Windows boot
- **Intuitive GUI**: User-friendly interface built with Fyne framework
- **Lightweight**: Minimal resource usage with native Windows API integration

## System Requirements

- **Operating System**: Windows 10/11 (64-bit)
- **Memory**: 50MB minimum RAM
- **Disk Space**: 20MB for installation

## Installation

### From Pre-built Binary

1. Download the latest release from the [Releases](https://github.com/Snax1210/gamma-changer/releases) page
2. Extract the archive to your desired location
3. Run `gammactl.exe` to launch the application

### Building from Source

1. Clone the repository:

```bash
git clone https://github.com/Snax1210/gamma-changer.git
cd gamma-changer
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

**Default build** (Ctrl+Alt+Key hotkeys):

```bash
go build -o gammactl.exe ./cmd/gammactl
```

**Single-key hotkey build** — triggers on a single key press, no modifier required:

```bash
go build -tags singlekey -o gammactl.exe ./cmd/gammactl
```

4. Run the application:

```bash
./gammactl.exe
```

## Usage

### Basic Usage

1. Launch `gammactl.exe`
2. Select the display you want to calibrate from the dropdown menu
3. Adjust the sliders for Gamma, Brightness, Contrast, and RGB Gain
4. Click "Apply" to apply the changes to your display
5. Save your settings as a preset for future use

### Command Line Interface

The application also provides a CLI for advanced users:

```bash
# Set gamma value
gammactl set --gamma 1.2

# Set brightness
gammactl set --brightness 0.3

# Apply a preset
gammactl preset apply "Night"

# Reset to default
gammactl reset

# Get current settings
gammactl get
```

## Parameters

### Gamma

- **Range**: 0.30 - 4.40
- **Default**: 1.0
- **Description**: Controls the overall gamma curve of the display. Higher values make the image brighter, lower values
  make it darker.

### Brightness

- **Range**: -1.00 - 1.00
- **Default**: 0.0
- **Description**: Adjusts the overall brightness level. Positive values increase brightness, negative values decrease
  it.

### Contrast

- **Range**: 0.10 - 3.00
- **Default**: 1.0
- **Description**: Controls the contrast ratio. Higher values increase contrast, lower values decrease it.

### RGB Gain

- **Range**: 0.0 - 2.0 (per channel)
- **Default**: 1.0 (all channels)
- **Description**: Independent adjustment for Red, Green, and Blue channels to fine-tune color balance.

## Preset Management

### Default Presets

The application comes with four built-in presets:

- **Default**: Standard display settings (Gamma: 1.0, Brightness: 0.0, Contrast: 1.0)
- **Office**: Optimized for office work (Gamma: 1.1, Brightness: 0.1, Contrast: 1.1)
- **Night**: Reduced brightness for nighttime use (Gamma: 0.9, Brightness: -0.3, Contrast: 0.9)
- **Coding**: Enhanced contrast for code readability (Gamma: 1.2, Brightness: 0.2, Contrast: 1.2)

### Creating Custom Presets

1. Adjust the display parameters to your desired values
2. Click the "Save Preset" button
3. Enter a name for your preset
4. Click "Save"

### Managing Presets

- **Apply**: Select a preset from the dropdown and click "Apply"
- **Rename**: Right-click on a preset and select "Rename"
- **Delete**: Right-click on a preset and select "Delete"

## Hotkey Configuration

Gamma Changer supports two hotkey modes, selected at build time via Go build tags:

### Hotkey Modes

| Mode                       | Build Command                             | Hotkey Format    | Use Case                                          |
|----------------------------|-------------------------------------------|------------------|---------------------------------------------------|
| **Default** (Ctrl+Alt+Key) | `go build ./cmd/gammactl`                 | `Ctrl+Alt+<Key>` | General purpose, avoids conflicts                 |
| **Single-key**             | `go build -tags singlekey ./cmd/gammactl` | `<Key>`          | Minimal input, e.g. dedicated keyboard, streaming |

### Setting Up Hotkeys

1. Click the "Bind" button next to the preset you want to assign a hotkey to
2. A dialog will appear showing the current hotkey mode
3. Press the desired key (A-Z or 0-9)
4. The hotkey is saved and applied immediately

### Hotkey Format

**Default mode** — `Ctrl+Alt+[Key]`:

- `Ctrl+Alt+1` through `Ctrl+Alt+9`: Number keys
- `Ctrl+Alt+A` through `Ctrl+Alt+Z`: Letter keys

**Single-key mode** — just the key itself:

- `1` through `9`: Number keys
- `A` through `Z`: Letter keys

## Configuration Files

### Main Configuration

- **Location**: `%USERCONFIGDIR%/gammactl/config.json`
- **Purpose**: Stores application settings, presets, and hotkey configurations

Example configuration:

```json
{
  "presets": {
    "Default": {
      "gamma": 1.0,
      "brightness": 0.0,
      "contrast": 1.0,
      "rgbGain": [
        1.0,
        1.0,
        1.0
      ]
    }
  },
  "hotkeys": {
    "Default": "Ctrl+Alt+1"
  },
  "autoStart": true
}
```

### Backup Files

- **Location**: `%USERCONFIGDIR%/gammactl/{display}_backup_ramp.json`
- **Purpose**: Stores original display settings for restoration
- **Format**: JSON file containing the original gamma ramp data

## Technical Architecture

### Project Structure

```
gamma-changer/
├── cmd/
│   └── gammactl/
│       ├── main.go              # GUI application entry point
│       ├── controller/          # MVC controllers
│       │   ├── state_controller.go
│       │   └── preset_controller.go
│       └── ui/                  # Fyne UI components
│           ├── app.go
│           ├── window.go
│           ├── sliders.go
│           ├── display_selector.go
│           ├── preset_list.go
│           ├── preset_dialog.go
│           └── hotkey_dialog.go
├── internal/
│   ├── app/
│   │   ├── app.go               # Core application logic
│   │   ├── config.go            # Configuration file management
│   │   ├── hotkeys.go           # Hotkey manager (runtime)
│   │   ├── hotkey_format_default.go  # Ctrl+Alt format (!singlekey)
│   │   ├── hotkey_format_single.go   # Single-key format (singlekey)
│   │   └── autostart_windows.go      # Windows auto-start
│   └── win/
│       └── gamma/
│           ├── ramp.go           # Gamma Ramp operations (Windows API)
│           └── display_enum.go   # Display enumeration
├── go.mod                       # Go module definition
└── go.sum                       # Dependency checksums
```

### Technology Stack

- **Language**: Go 1.24.10
- **GUI Framework**: Fyne v2.7.2
- **Hotkey Library**: golang.design/x/hotkey v0.4.1
- **System Calls**: golang.org/x/sys v0.40.0

### Core Components

1. **Gamma Ramp Manager**: Handles Windows API calls for gamma ramp manipulation
2. **Display Enumerator**: Discovers and manages multiple displays
3. **Preset Manager**: Saves, loads, and applies display presets
4. **Hotkey Manager**: Registers and handles global hotkey events
5. **Configuration Manager**: Manages application settings and persistence

### Windows API Integration

The application uses the following Windows APIs:

- `EnumDisplayMonitors`: Enumerate all connected displays
- `GetDeviceGammaRamp`: Retrieve current gamma ramp
- `SetDeviceGammaRamp`: Apply new gamma ramp values
- Registry APIs: Manage auto-start functionality

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards and best practices
- Write clear, documented code
- Test your changes thoroughly
- Update documentation as needed

## Support

If you encounter any issues or have questions:

- Open an issue on [GitHub Issues](https://github.com/Snax1210/gamma-changer/issues)
- Check existing documentation and discussions
- Contact the maintainers

## Acknowledgments

- Built with [Fyne](https://fyne.io/) - A cross-platform GUI toolkit for Go
- Hotkey functionality powered by [golang.design/x/hotkey](https://github.com/golang-design/hotkey)

## Changelog

### Version 1.1.0

- Added single-key hotkey mode (`-tags singlekey`) alongside the default Ctrl+Alt+Key mode
- Hotkey format is now configurable via Go build tags at compile time

### Version 1.0.1

- Optimize code structure

### Version 1.0.0

- Initial release
- Real-time display parameter adjustment
- Preset management system
- Global hotkey support (Ctrl+Alt+Key)
- Multi-monitor support
- Auto-start functionality

---

**Note**: This software modifies display settings at the system level. Always backup your original settings before
making changes. The authors are not responsible for any damage to your display hardware.
