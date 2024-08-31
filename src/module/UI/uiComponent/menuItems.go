package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
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

	t.setHotKey()
	language.Center.RegisterChange(t.changeLanguageHandler)
}

// 设置热键
func (t *MenuItemT) setHotKey() {
	_ = eventCenter.Event.Publish(events.SetHotKey, events.SetHotKeyData{
		Key: t.hKList[0],
	})
	_ = eventCenter.Event.Publish(events.SetHotKey, events.SetHotKeyData{
		Key: t.hKList[1],
	})
	_ = eventCenter.Event.Publish(events.SetHotKey, events.SetHotKeyData{
		Key: t.hKList[2],
	})
	_ = eventCenter.Event.Publish(events.SetHotKey, events.SetHotKeyData{
		Key: t.hKList[3],
	})

	//t.ChangeLanguage(t.languageTyp, false)
}

// 设置热键弹窗
func (t *MenuItemT) setHotKeyPop() {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var tmpList = t.hKList

	cmd, _ := Dialog{AssignTo: &dlg, Title: language.Center.Get(language.SetHotKeyWindowTitleStr),
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			Label{Text: language.Center.Get(language.RecordStr), ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[0], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[0], OnCurrentIndexChanged: func() {
				tmpList[0] = t.hKBox[0].Text()
			}},

			Label{Text: language.Center.Get(language.PlaybackStr), ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[1], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[1], OnCurrentIndexChanged: func() {
				tmpList[1] = t.hKBox[1].Text()
			}},

			Label{Text: language.Center.Get(language.PauseStr), ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[2], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[2], OnCurrentIndexChanged: func() {
				tmpList[2] = t.hKBox[2].Text()
			}},

			Label{Text: language.Center.Get(language.StopStr), ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[3], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[3], OnCurrentIndexChanged: func() {
				tmpList[3] = t.hKBox[3].Text()
			}},

			PushButton{ColumnSpan: 4, Text: language.Center.Get(language.ResetStr), OnClicked: func() {
				tmpList = [4]string{"F7", "F8", "F9", "F10"}
				_ = t.hKBox[0].SetText(tmpList[0])
				_ = t.hKBox[1].SetText(tmpList[1])
				_ = t.hKBox[2].SetText(tmpList[2])
				_ = t.hKBox[3].SetText(tmpList[3])
			}},

			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.Center.Get(language.OKStr), OnClicked: func() {
				M := make(map[string]bool)
				for _, v := range tmpList {
					if M[v] {
						walk.MsgBox(dlg, language.Center.Get(language.ErrWindowTitleStr), language.Center.Get(language.SetHotKeyErrMessageStr), walk.MsgBoxIconInformation)
						return
					} else {
						M[v] = true
					}
				}
				t.hKList = tmpList
				dlg.Accept()
			}},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.Center.Get(language.CancelStr), OnClicked: func() { dlg.Cancel() }},
		},
	}.Run(t.mw)

	if cmd == walk.DlgCmdOK {
		t.setHotKey()
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
