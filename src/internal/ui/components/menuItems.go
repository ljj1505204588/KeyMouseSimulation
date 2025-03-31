package uiComponent

import (
	gene "KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/share/enum"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	component "KeyMouseSimulation/module/baseComponent"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sort"
	"sync"
	"time"
)

// MenuItemT 设置栏
type MenuItemT struct {
	mw **walk.MainWindow
	sync.Once

	//工具
	settingMenu        *walk.Action
	setHotkeyAction    *walk.Action
	acceptPB, cancelPB *walk.PushButton

	languageMenu *walk.Action
	helpMenu     *walk.Action
	aboutAction  *walk.Action

	menuItems []MenuItem
}

func (t *MenuItemT) MenuItems(mw **walk.MainWindow) []MenuItem {
	t.mw = mw
	t.Once.Do(t.Init)

	return t.menuItems
}

func (t *MenuItemT) Init() {
	t.menuItems = []MenuItem{
		Menu{AssignActionTo: &t.settingMenu, Items: []MenuItem{
			Action{AssignTo: &t.setHotkeyAction, OnTriggered: t.setHotKeyPop},
			Menu{AssignActionTo: &t.languageMenu, Items: []MenuItem{
				Action{Text: string(component.English), OnTriggered: func() {
					component.Center.SetLanguage(component.English)
				}},
				Action{Text: string(component.Chinese), OnTriggered: func() {
					component.Center.SetLanguage(component.Chinese)
				}},
			},
			},
		}},
		Menu{AssignActionTo: &t.helpMenu, Items: []MenuItem{
			Action{AssignTo: &t.aboutAction, OnTriggered: t.showAboutBoxAction},
		}},
	}

	component.Center.RegisterChange(t.changeLanguageHandler)
}

// 设置热键弹窗
func (t *MenuItemT) setHotKeyPop() {
	type hkSetT struct {
		name   enum.HotKey
		hk     component.HotKeyI
		box    *walk.ComboBox
		widget []Widget
	}

	var setM = make(map[component.HotKeyI]string)
	var resetMethod []func()

	var hkSets []hkSetT
	for name, hk := range component.GetAllHk() {
		var set = hkSetT{name: name, hk: hk, box: &walk.ComboBox{}}
		set.widget = []Widget{
			Label{Text: component.Center.Get(name.Language()), ColumnSpan: 1},
			ComboBox{AssignTo: &set.box, ColumnSpan: 1, Model: keyMouTool.VKCodeStringKeys, Editable: true, Value: hk.Key(), OnCurrentIndexChanged: func() {
				setM[set.hk] = set.box.Text()
			}},
		}

		hkSets = append(hkSets, set)

		var defKey = hk.DefaultKey()
		resetMethod = append(resetMethod, func() {
			_ = set.box.SetText(defKey)
			setM[hk] = defKey
		})
		setM[hk] = hk.Key() // 校验重复性
	}
	sort.Slice(hkSets, func(i, j int) bool {
		return hkSets[i].name < hkSets[j].name
	})

	var hkWidget []Widget
	for _, set := range hkSets {
		hkWidget = append(hkWidget, set.widget...)
	}

	var dlg *walk.Dialog
	hkWidget = append(hkWidget, []Widget{
		PushButton{ColumnSpan: 4, Text: component.Center.Get(component.ResetStr), OnClicked: func() {
			for _, reset := range resetMethod {
				reset()
			}
		}},

		PushButton{AssignTo: &t.acceptPB, ColumnSpan: 2, Text: component.Center.Get(component.OKStr), OnClicked: func() {
			var keys []string
			for _, key := range setM {
				keys = append(keys, key)
			}

			if len(gene.RemoveDuplicate(keys)) != len(keys) {
				walk.MsgBox(dlg, component.Center.Get(component.ErrWindowTitleStr), component.Center.Get(component.SetHotKeyErrMessageStr), walk.MsgBoxIconInformation)
				return
			}

			dlg.Accept()
		}},
		PushButton{AssignTo: &t.cancelPB, ColumnSpan: 2, Text: component.Center.Get(component.CancelStr), OnClicked: func() { dlg.Cancel() }},
	}...)

	cmd, _ := Dialog{AssignTo: &dlg, Title: component.Center.Get(component.SetHotKeyWindowTitleStr),
		DefaultButton: &t.acceptPB, CancelButton: &t.cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: hkWidget,
	}.Run(*t.mw)

	if cmd == walk.DlgCmdOK {
		tryPublishErr(component.MulSetKey(setM))
		component.Center.Refresh()
	}
}

// 初始化校验
func (t *MenuItemT) initCheck() bool {
	for _, per := range []*walk.Action{
		t.settingMenu,
		t.setHotkeyAction,
		t.languageMenu,
		t.helpMenu,
		t.aboutAction,
	} {
		if per == nil {
			return false
		}
	}
	return true
}

// 修改语言
func (t *MenuItemT) changeLanguageHandler() {
	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	tryPublishErr(t.settingMenu.SetText(component.Center.Get(component.MenuSettingStr)))
	tryPublishErr(t.languageMenu.SetText(component.Center.Get(component.MenuItemLanguageStr)))
	tryPublishErr(t.setHotkeyAction.SetText(component.Center.Get(component.ActionSetHotKeyStr)))
	tryPublishErr(t.helpMenu.SetText(component.Center.Get(component.MenuHelpStr)))
	tryPublishErr(t.aboutAction.SetText(component.Center.Get(component.ActionAboutStr)))
}

// 系统信息
func (t *MenuItemT) showAboutBoxAction() {
	walk.MsgBox(*t.mw, component.Center.Get(component.AboutWindowTitleStr), component.Center.Get(component.AboutMessageStr), walk.MsgBoxIconInformation)
}
