package uiComponent

import (
	gene "KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/share/enum"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/module/language"
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
				Action{Text: string(language.English), OnTriggered: func() {
					language.Center.SetLanguage(language.English)
				}},
				Action{Text: string(language.Chinese), OnTriggered: func() {
					language.Center.SetLanguage(language.Chinese)
				}},
			},
			},
		}},
		Menu{AssignActionTo: &t.helpMenu, Items: []MenuItem{
			Action{AssignTo: &t.aboutAction, OnTriggered: t.showAboutBoxAction},
		}},
	}

	language.Center.RegisterChange(t.changeLanguageHandler)
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
			Label{Text: language.Center.Get(name.Language()), ColumnSpan: 1},
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
		PushButton{ColumnSpan: 4, Text: language.Center.Get(language.ResetStr), OnClicked: func() {
			for _, reset := range resetMethod {
				reset()
			}
		}},

		PushButton{AssignTo: &t.acceptPB, ColumnSpan: 2, Text: language.Center.Get(language.OKStr), OnClicked: func() {
			var keys []string
			for _, key := range setM {
				keys = append(keys, key)
			}

			if len(gene.RemoveDuplicate(keys)) != len(keys) {
				walk.MsgBox(dlg, language.Center.Get(language.ErrWindowTitleStr), language.Center.Get(language.SetHotKeyErrMessageStr), walk.MsgBoxIconInformation)
				return
			}

			dlg.Accept()
		}},
		PushButton{AssignTo: &t.cancelPB, ColumnSpan: 2, Text: language.Center.Get(language.CancelStr), OnClicked: func() { dlg.Cancel() }},
	}...)

	cmd, _ := Dialog{AssignTo: &dlg, Title: language.Center.Get(language.SetHotKeyWindowTitleStr),
		DefaultButton: &t.acceptPB, CancelButton: &t.cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: hkWidget,
	}.Run(*t.mw)

	if cmd == walk.DlgCmdOK {
		tryPublishErr(component.MulSetKey(setM))
		language.Center.Refresh()
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

	tryPublishErr(t.settingMenu.SetText(language.Center.Get(language.MenuSettingStr)))
	tryPublishErr(t.languageMenu.SetText(language.Center.Get(language.MenuItemLanguageStr)))
	tryPublishErr(t.setHotkeyAction.SetText(language.Center.Get(language.ActionSetHotKeyStr)))
	tryPublishErr(t.helpMenu.SetText(language.Center.Get(language.MenuHelpStr)))
	tryPublishErr(t.aboutAction.SetText(language.Center.Get(language.ActionAboutStr)))
}

// 系统信息
func (t *MenuItemT) showAboutBoxAction() {
	walk.MsgBox(*t.mw, language.Center.Get(language.AboutWindowTitleStr), language.Center.Get(language.AboutMessageStr), walk.MsgBoxIconInformation)
}
