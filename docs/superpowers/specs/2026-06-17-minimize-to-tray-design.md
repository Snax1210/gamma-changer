# 关闭窗口最小化到系统托盘 — 设计文档

日期：2026-06-17

## 目标

点击主窗口关闭按钮（X）时不退出程序，而是隐藏窗口到系统托盘。托盘提供菜单用于恢复窗口或真正退出程序。

## 背景

- 应用基于 Fyne v2.7.2，目标平台 Windows（`//go:build windows`）。
- 入口 `cmd/gammactl/main.go`，主窗口通过 `w.ShowAndRun()` 启动。
- `fyne.io/systray` 已作为间接依赖存在，Fyne 原生支持系统托盘。

## 关键 API（Fyne 2.7）

`desktop.App` 接口提供 `SetSystemTrayWindow(fyne.Window)`（自 2.7）：在 Windows 上**左键点击托盘图标即显示窗口**，并自动设置 `w.SetCloseIntercept(w.Hide)`（关闭→隐藏）。若同时设置了菜单，菜单改为**右键**弹出。此外 Fyne 会通过 `addMissingQuitForMenu` 自动为托盘菜单追加 "Quit" 项，无需手动添加退出项。菜单文案沿用项目现有英文风格（"Show Window"）。

## 文件改动（均限定 `//go:build windows`）

1. 新增图标资源 `cmd/gammactl/ui/assets/tray.png`（从用户提供的 PNG 拷贝而来，embed 要求文件在模块内）。

2. 新增 `cmd/gammactl/ui/tray.go`：
   - `//go:embed assets/tray.png` 嵌入图标，封装为 `fyne.StaticResource`。
   - `func SetupSystemTray(gui fyne.App, w fyne.Window) bool`：
     - `desk, ok := gui.(desktop.App)`，失败返回 `false`（优雅降级）。
     - 构建菜单含 "Show Window"(`w.Show`)，Quit 由 Fyne 自动追加。
     - 依次 `SetSystemTrayIcon` → `SetSystemTrayMenu` → `SetSystemTrayWindow`，返回 `true`。

3. 修改 `cmd/gammactl/main.go`：在 `w.ShowAndRun()` 之前调用 `ui.SetupSystemTray(gui, w)`。关闭拦截由 `SetSystemTrayWindow` 内部设置，仅在托盘可用时生效。

## 数据流 / 行为

- 启动 → 正常显示窗口。
- 点 X → 拦截 → `w.Hide()`，进程驻留托盘。
- 左键点击托盘图标 → 显示窗口。
- 右键托盘菜单 "Show Window" → 显示窗口；"Quit" → 退出程序。

## 错误处理

- `desktop.App` 断言失败 → 不设托盘、不拦截关闭，保持默认（关闭即退出）行为。

## 验证

- `go build ./...` 通过、`go vet ./...` 无误。
- 手动验证：运行后点关闭 → 窗口消失但进程驻留、托盘出现图标 → 左键点击恢复、右键菜单可恢复/退出。