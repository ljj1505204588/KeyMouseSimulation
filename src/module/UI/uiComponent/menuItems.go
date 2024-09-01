package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	gene "KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/share/events"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"time"
)

// MenuItemT 设置栏
type MenuItemT struct {
	mw *walk.MainWindow
	sync.Once

	//工具
	settingMenu     *walk.Action
	setHotkeyAction *walk.Action

	languageMenu *walk.Action
	helpMenu     *walk.Action
	aboutAction  *walk.Action

	menuItems []MenuItem
}

func (t *MenuItemT) MenuItems(mw *walk.MainWindow) []MenuItem {
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
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var hkWidget []Widget

	var resetMethod []func()
	var setM = make(map[component.HotKeyI]string)
	for name, hk := range component.GetAllHk() {
		var box = &walk.ComboBox{}
		hkWidget = append(hkWidget, []Widget{
			Label{Text: language.Center.Get(name.Language()), ColumnSpan: 1},
			ComboBox{AssignTo: &box, ColumnSpan: 1, Model: keyMouTool.VKCodeStringKeys, Editable: true, Value: hk.Key(), OnCurrentIndexChanged: func() {
				setM[hk] = box.Text()
			}},
		}...)

		resetMethod = append(resetMethod, func() {
			_ = box.SetText(hk.DefaultKey())
			setM[hk] = hk.DefaultKey()
		})
		setM[hk] = hk.Key() // 校验重复性
	}

	cmd, _ := Dialog{AssignTo: &dlg, Title: language.Center.Get(language.SetHotKeyWindowTitleStr),
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: append(hkWidget, []Widget{

			PushButton{ColumnSpan: 4, Text: language.Center.Get(language.ResetStr), OnClicked: func() {
				for _, reset := range resetMethod {
					reset()
				}
			}},

			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.Center.Get(language.OKStr), OnClicked: func() {
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
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.Center.Get(language.CancelStr), OnClicked: func() { dlg.Cancel() }},
		}...),
	}.Run(t.mw)

	if cmd == walk.DlgCmdOK {
		if err := component.MulSetKey(setM); err != nil {
			_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
		}
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

	_ = t.settingMenu.SetText(language.Center.Get(language.MenuSettingStr))
	_ = t.languageMenu.SetText(language.Center.Get(language.MenuItemLanguageStr))
	_ = t.setHotkeyAction.SetText(language.Center.Get(language.ActionSetHotKeyStr))
	_ = t.helpMenu.SetText(language.Center.Get(language.MenuHelpStr))
	_ = t.aboutAction.SetText(language.Center.Get(language.ActionAboutStr))

}

// 系统信息
func (t *MenuItemT) showAboutBoxAction() {
	walk.MsgBox(t.mw, language.Center.Get(language.AboutWindowTitleStr), language.Center.Get(language.AboutMessageStr), walk.MsgBoxIconInformation)
}
