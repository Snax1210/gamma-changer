//go:build windows

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gamma-changer/internal/app"
	"gamma-changer/internal/win/gamma"

	"fyne.io/fyne/v2"
	fapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	st := app.LoadStateOrDefault()

	// ✅ core 需要可重建（避免 State() 拷贝/缓存导致 UI 不更新）
	var core *app.App
	core = app.New(st)

	// Hotkey manager：启动即注册
	hkMgr := app.NewHotkeyManager()
	_ = hkMgr.ApplyFromState(core.State(), core)

	// displays
	ds, err := gamma.ListDisplays()
	if err == nil && st.SelectedDisplay == "" && len(ds) > 0 {
		core.SetSelectedDisplay(ds[0].Name)
	}
	displayNames := make([]string, 0, len(ds))
	for _, d := range ds {
		displayNames = append(displayNames, d.Name)
	}

	gui := fapp.New()
	w := gui.NewWindow("Gamma Changer")
	w.Resize(fyne.NewSize(780, 640))

	status := widget.NewLabel("Ready")

	// ---- debounce apply ----
	var t *time.Timer
	const debounce = 80 * time.Millisecond

	// ✅ 统一：从磁盘重载 core + 重载热键（确保 UI/逻辑一致）
	reloadCoreFromDisk := func() {
		st2 := app.LoadStateOrDefault()
		core = app.New(st2)
		_ = hkMgr.ApplyFromState(core.State(), core)
	}

	applyNow := func() {
		if err := core.ApplyCurrent(); err != nil {
			status.SetText("Apply failed: " + err.Error())
			return
		}
		s := core.State()
		status.SetText(fmt.Sprintf("OK  γ=%.2f  b=%.2f  c=%.2f  (%s)",
			s.Params.Gamma, s.Params.Brightness, s.Params.Contrast, s.SelectedDisplay))
		_ = app.SaveState(s)
	}

	scheduleApply := func() {
		if t != nil {
			t.Stop()
		}
		t = time.AfterFunc(debounce, func() {
			applyNow()
		})
	}

	// Display dropdown
	displaySelect := widget.NewSelect(displayNames, func(s string) {
		core.SetSelectedDisplay(s)
		scheduleApply()
	})
	displaySelect.PlaceHolder = "Select display"
	if core.State().SelectedDisplay != "" {
		displaySelect.SetSelected(core.State().SelectedDisplay)
	}

	// Sliders + labels
	row := func(title string, slider *widget.Slider, value *widget.Label) fyne.CanvasObject {
		return container.NewBorder(nil, nil, widget.NewLabel(title), value, slider)
	}

	// Gamma
	gammaSlider = widget.NewSlider(0.30, 4.40)
	gammaSlider.Step = 0.01
	gammaSlider.Value = core.State().Params.Gamma
	gammaValue := widget.NewLabel(fmt.Sprintf("%.2f", gammaSlider.Value))
	gammaSlider.OnChanged = func(v float64) {
		core.UpdateParams(func(p *app.Params) { p.Gamma = v })
		gammaValue.SetText(fmt.Sprintf("%.2f", v))
		scheduleApply()
	}

	// Brightness
	brightSlider = widget.NewSlider(-1.00, 1.00)
	brightSlider.Step = 0.01
	brightSlider.Value = core.State().Params.Brightness
	brightValue := widget.NewLabel(fmt.Sprintf("%.2f", brightSlider.Value))
	brightSlider.OnChanged = func(v float64) {
		core.UpdateParams(func(p *app.Params) { p.Brightness = v })
		brightValue.SetText(fmt.Sprintf("%.2f", v))
		scheduleApply()
	}

	// Contrast
	contrastSlider = widget.NewSlider(0.10, 3.00)
	contrastSlider.Step = 0.01
	contrastSlider.Value = core.State().Params.Contrast
	contrastValue := widget.NewLabel(fmt.Sprintf("%.2f", contrastSlider.Value))
	contrastSlider.OnChanged = func(v float64) {
		core.UpdateParams(func(p *app.Params) { p.Contrast = v })
		contrastValue.SetText(fmt.Sprintf("%.2f", v))
		scheduleApply()
	}

	// Buttons
	btnReset := widget.NewButton("Reset (Restore backup)", func() {
		if err := core.Reset(); err != nil {
			status.SetText("Reset failed: " + err.Error())
			return
		}
		status.SetText("Reset OK")
		_ = app.SaveState(core.State())
	})

	btnApply := widget.NewButton("Apply Now", applyNow)

	// ======== Preset Hotkey Binding & Editing UI ========

	// 删除 preset 后：preset index 会变，需要重排 hotkeys（把 >idx 的 key -1）
	shiftHotkeysAfterDelete := func(hk map[string]string, deletedIdx int) map[string]string {
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

	bindingHint := widget.NewLabel("Hotkeys: click Bind then press a key (Ctrl+Alt+<Key>)")
	bindingHint.Wrapping = fyne.TextWrapWord

	presetList := container.NewVBox()

	// 刷新预设列表：✅ 每次都从磁盘加载，保证新增预设立即显示
	var refreshPresetList func()

	// 绑定快捷键（弹窗）：✅ completed 防止 Hide() 触发 canceled 覆盖
	bindDialog := func(idx int, name string) {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Press A-Z or 0-9 (no need to hold Ctrl/Alt)")
		entry.Wrapping = fyne.TextWrapOff

		info := widget.NewLabel(fmt.Sprintf(
			"Binding preset '%s'\nResult will be: Ctrl+Alt+<Key>\n(Press one key: A-Z / 0-9)",
			name,
		))
		info.Wrapping = fyne.TextWrapWord

		var d dialog.Dialog
		completed := false

		entry.OnChanged = func(s string) {
			if s == "" || completed {
				return
			}
			r := []rune(s)[0]
			key := strings.ToUpper(string(r))

			if !(len(key) == 1 && ((key[0] >= '0' && key[0] <= '9') || (key[0] >= 'A' && key[0] <= 'Z'))) {
				status.SetText("Bind: only A-Z / 0-9 supported.")
				entry.SetText("")
				return
			}

			spec := "Ctrl+Alt+" + key

			stDisk := app.LoadStateOrDefault()
			if stDisk.Hotkeys == nil {
				stDisk.Hotkeys = map[string]string{}
			}
			stDisk.Hotkeys[strconv.Itoa(idx)] = spec

			_ = app.SaveState(stDisk)

			// ✅ 重载 core/hotkeys + 刷新 UI
			reloadCoreFromDisk()

			completed = true
			status.SetText(fmt.Sprintf("Bound preset '%s' -> %s", name, spec))
			bindingHint.SetText("Hotkeys: click Bind then press a key (Ctrl+Alt+<Key>)")

			if d != nil {
				d.Hide()
			}
			refreshPresetList()
		}

		content := container.NewVBox(info, entry)

		d = dialog.NewCustomConfirm(
			"Bind Hotkey",
			"Cancel",
			"",
			content,
			func(_ bool) {
				if completed {
					return
				}
				status.SetText("Bind canceled.")
				bindingHint.SetText("Hotkeys: click Bind then press a key (Ctrl+Alt+<Key>)")
			},
			w,
		)

		bindingHint.SetText(fmt.Sprintf("Binding preset '%s'... (via dialog)", name))
		d.Show()

		time.AfterFunc(10*time.Millisecond, func() {
			w.Canvas().Focus(entry)
		})
	}

	refreshPresetList = func() {
		presetList.Objects = nil

		// ✅ 每次都从磁盘读，保证新增/改名/删除立刻体现在 UI
		stDisk := app.LoadStateOrDefault()

		hotkeyLabel := func(idx int) string {
			if stDisk.Hotkeys == nil {
				return "-"
			}
			if v, ok := stDisk.Hotkeys[strconv.Itoa(idx)]; ok && v != "" {
				return v
			}
			return "-"
		}

		for i, p := range stDisk.Presets {
			idx := i
			name := p.Name

			lblName := widget.NewLabel(name)
			lblHK := widget.NewLabel(hotkeyLabel(idx))

			btnPresetApply := widget.NewButton("Apply", func() {
				// ✅ Apply 之前重载 core，确保 core 的 presets 跟磁盘一致
				reloadCoreFromDisk()

				if err := core.ApplyPreset(idx); err != nil {
					status.SetText("Preset failed: " + err.Error())
					return
				}

				// sync sliders
				s2 := core.State()
				gammaSlider.SetValue(s2.Params.Gamma)
				brightSlider.SetValue(s2.Params.Brightness)
				contrastSlider.SetValue(s2.Params.Contrast)

				status.SetText("Preset applied: " + name)
				_ = app.SaveState(s2)
			})

			btnBind := widget.NewButton("Bind", func() {
				bindDialog(idx, name)
			})

			btnClear := widget.NewButton("Clear", func() {
				st2 := app.LoadStateOrDefault()
				if st2.Hotkeys == nil {
					st2.Hotkeys = map[string]string{}
				}
				delete(st2.Hotkeys, strconv.Itoa(idx))
				_ = app.SaveState(st2)

				reloadCoreFromDisk()
				status.SetText("Cleared hotkey for: " + name)
				refreshPresetList()
			})

			// 保存当前参数到该预设
			btnSave := widget.NewButton("Save", func() {
				// 用当前 core 的参数保存进磁盘预设
				cur := core.State().Params

				st2 := app.LoadStateOrDefault()
				if idx < 0 || idx >= len(st2.Presets) {
					return
				}
				st2.Presets[idx].Params = cur
				_ = app.SaveState(st2)

				reloadCoreFromDisk()
				status.SetText("Preset saved from current: " + name)
				refreshPresetList()
			})

			// 重命名
			btnRename := widget.NewButton("Rename", func() {
				entry := widget.NewEntry()
				entry.SetText(name)
				d := dialog.NewForm("Rename Preset", "OK", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("Name", entry),
					},
					func(ok bool) {
						if !ok {
							return
						}
						newName := strings.TrimSpace(entry.Text)
						if newName == "" {
							status.SetText("Preset name cannot be empty.")
							return
						}

						st2 := app.LoadStateOrDefault()
						if idx < 0 || idx >= len(st2.Presets) {
							return
						}
						st2.Presets[idx].Name = newName
						_ = app.SaveState(st2)

						reloadCoreFromDisk()
						status.SetText("Preset renamed.")
						refreshPresetList()
					}, w)
				d.Resize(fyne.NewSize(420, 160))
				d.Show()
			})

			// 删除（并重排 hotkeys）
			btnDelete := widget.NewButton("Delete", func() {
				d := dialog.NewConfirm("Delete Preset",
					fmt.Sprintf("Delete preset '%s' ?", name),
					func(ok bool) {
						if !ok {
							return
						}
						st2 := app.LoadStateOrDefault()
						if idx < 0 || idx >= len(st2.Presets) {
							return
						}
						st2.Presets = append(st2.Presets[:idx], st2.Presets[idx+1:]...)
						st2.Hotkeys = shiftHotkeysAfterDelete(st2.Hotkeys, idx)

						_ = app.SaveState(st2)

						reloadCoreFromDisk()
						status.SetText("Preset deleted: " + name)
						refreshPresetList()
					}, w)
				d.Show()
			})

			line := container.NewGridWithColumns(8,
				lblName,
				lblHK,
				btnPresetApply,
				btnBind,
				btnClear,
				btnSave,
				btnRename,
				btnDelete,
			)
			presetList.Add(line)
		}

		presetList.Refresh()
	}

	// 新建预设（从当前参数复制）
	btnNewPreset := widget.NewButton("New Preset (from current)", func() {
		entry := widget.NewEntry()
		entry.SetText("New Preset")
		d := dialog.NewForm("New Preset", "Create", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Name", entry),
			},
			func(ok bool) {
				if !ok {
					return
				}
				name := strings.TrimSpace(entry.Text)
				if name == "" {
					status.SetText("Preset name cannot be empty.")
					return
				}

				// ✅ 写磁盘：Presets append
				st2 := app.LoadStateOrDefault()
				st2.Presets = append(st2.Presets, app.Preset{
					Name:   name,
					Params: core.State().Params, // 用当前参数作为模板
				})
				if st2.Hotkeys == nil {
					st2.Hotkeys = map[string]string{}
				}
				_ = app.SaveState(st2)

				// ✅ 立即生效：重载 + 刷新
				reloadCoreFromDisk()
				status.SetText("Preset created: " + name)
				refreshPresetList()
			}, w)
		d.Resize(fyne.NewSize(420, 160))
		d.Show()
	})

	// 初次渲染列表
	refreshPresetList()

	// 主布局
	content := container.NewVBox(
		widget.NewLabelWithStyle("Gamma / Brightness / Contrast (Software LUT, Real-time)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Display"),
		displaySelect,

		widget.NewSeparator(),
		widget.NewLabelWithStyle("Presets & Hotkeys (Global)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnNewPreset),
		bindingHint,
		presetList,

		widget.NewSeparator(),
		row("Gamma", gammaSlider, gammaValue),
		row("Brightness", brightSlider, brightValue),
		row("Contrast", contrastSlider, contrastValue),

		container.NewHBox(btnApply, btnReset),
		widget.NewSeparator(),
		status,
	)

	w.SetContent(content)

	// 启动时应用一次（可选）
	scheduleApply()

	w.ShowAndRun()
}

// sliders declared here so preset can sync them
var (
	gammaSlider    *widget.Slider
	brightSlider   *widget.Slider
	contrastSlider *widget.Slider
)
