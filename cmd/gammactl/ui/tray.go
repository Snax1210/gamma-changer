//go:build windows

package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

//go:embed assets/tray.png
var trayIconData []byte

// trayIcon 是嵌入的系统托盘图标资源。
var trayIcon = fyne.NewStaticResource("tray.png", trayIconData)

// SetupSystemTray 配置系统托盘。
//
// 关闭主窗口时隐藏到托盘而非退出，左键点击托盘图标恢复窗口，
// 右键菜单提供 "Show Window" 与自动追加的 "Quit"。
// 当运行环境不支持桌面托盘时返回 false 且不做任何改动，
// 此时窗口保持默认的“关闭即退出”行为。
func SetupSystemTray(gui fyne.App, w fyne.Window) bool {
	desk, ok := gui.(desktop.App)
	if !ok {
		return false
	}

	menu := fyne.NewMenu("Gamma Changer",
		fyne.NewMenuItem("Show Window", w.Show),
	)

	desk.SetSystemTrayIcon(trayIcon)
	desk.SetSystemTrayMenu(menu)
	// SetSystemTrayWindow 同时设置“关闭隐藏到托盘”和“左键点击恢复窗口”。
	desk.SetSystemTrayWindow(w)
	return true
}
