package BaseComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module2/language"
	"KeyMouseSimulation/share/events"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

// SystemT 系统按钮
type SystemT struct {
	*BaseT

	//系统状态
	statusLabel *walk.Label
	statusEdit  *walk.LineEdit
	errorLabel  *walk.Label
	errorEdit   *walk.TextEdit

	//工具
	settingMenu     *walk.Action
	setHotkeyAction *walk.Action

	languageMenu *walk.Action
	helpMenu     *walk.Action
	aboutAction  *walk.Action

	widgets   []Widget
	menuItems []MenuItem
}

func (t *SystemT) Init(base *BaseT) {
	t.BaseT = base
	t.BaseT.registerChangeLanguage(t.changeLanguageHandler)

	t.widgets = []Widget{
		//当前状态
		Label{AssignTo: &t.statusLabel, ColumnSpan: 2},
		LineEdit{AssignTo: &t.statusEdit, ColumnSpan: 6, ReadOnly: true},

		//错误信息
		Label{AssignTo: &t.errorLabel, ColumnSpan: 2},
		TextEdit{AssignTo: &t.errorEdit, ColumnSpan: 6, ReadOnly: true},
	}
	t.menuItems = []MenuItem{
		Menu{AssignActionTo: &t.settingMenu, Items: []MenuItem{
			Action{AssignTo: &t.setHotkeyAction, OnTriggered: t.setHotKey},
			Menu{AssignActionTo: &t.languageMenu, Items: []MenuItem{
				Action{Text: string(language.English), OnTriggered: func() {
					t.ChangeLanguage(language.English, false)
				}},
				Action{Text: string(language.Chinese), OnTriggered: func() {
					t.ChangeLanguage(language.Chinese, false)
				}},
			},
			},
		}},
		Menu{AssignActionTo: &t.helpMenu, Items: []MenuItem{
			Action{AssignTo: &t.aboutAction, OnTriggered: t.showAboutBoxAction},
		}},
	}

	eventCenter.Event.Register(events.ServerError, t.subShowError)
	eventCenter.Event.Register(events.ServerChange, t.subServerChange)
}

func (t *SystemT) DisPlay() []Widget {
	return t.widgets
}

func (t *SystemT) MenuItems() []MenuItem {
	return t.menuItems
}

func (t *SystemT) MWPoint() **walk.MainWindow {
	return &t.mw
}

// --------------------------------------- 基础功能 ----------------------------------------------

// 初始化校验
func (t *SystemT) initCheck() bool {
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
	return t.statusLabel != nil && t.statusEdit != nil && t.errorLabel != nil && t.errorEdit != nil
}

// 修改语言
func (t *SystemT) changeLanguageHandler(typ language.LanguageTyp) {
	var m = t.languageMap

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	_ = t.settingMenu.SetText(m[language.MenuSettingStr])
	_ = t.languageMenu.SetText(m[language.MenuItemLanguageStr])
	_ = t.setHotkeyAction.SetText(m[language.ActionSetHotKeyStr])
	_ = t.helpMenu.SetText(m[language.MenuHelpStr])
	_ = t.aboutAction.SetText(m[language.ActionAboutStr])
	_ = t.statusLabel.SetText(m[language.StatusLabelStr])
	_ = t.errorLabel.SetText(m[language.ErrorLabelStr])

}

// 系统信息
func (t *SystemT) showAboutBoxAction() {
	walk.MsgBox(t.mw, t.languageMap[language.AboutWindowTitleStr], t.languageMap[language.AboutMessageStr], walk.MsgBoxIconInformation)
}

// 设置热键
func (t *SystemT) setHotKey() {
	if cmd, _ := t.setHotKeyPop(); cmd == walk.DlgCmdOK {
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

		t.ChangeLanguage(t.languageTyp, false)
	}
}

// 设置热键弹窗
func (t *SystemT) setHotKeyPop() (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var tmpList = t.hKList

	cmd, err := Dialog{AssignTo: &dlg, Title: t.languageMap[language.SetHotKeyWindowTitleStr],
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			Label{Text: t.languageMap[language.RecordStr], ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[0], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[0], OnCurrentIndexChanged: func() {
				tmpList[0] = t.hKBox[0].Text()
			}},

			Label{Text: t.languageMap[language.PlaybackStr], ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[1], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[1], OnCurrentIndexChanged: func() {
				tmpList[1] = t.hKBox[1].Text()
			}},

			Label{Text: t.languageMap[language.PauseStr], ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[2], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[2], OnCurrentIndexChanged: func() {
				tmpList[2] = t.hKBox[2].Text()
			}},

			Label{Text: t.languageMap[language.StopStr], ColumnSpan: 1},
			ComboBox{AssignTo: &t.hKBox[3], ColumnSpan: 1, Model: t.keyList, Editable: true, Value: t.hKList[3], OnCurrentIndexChanged: func() {
				tmpList[3] = t.hKBox[3].Text()
			}},

			PushButton{ColumnSpan: 4, Text: t.languageMap[language.ResetStr], OnClicked: func() {
				tmpList = [4]string{"F7", "F8", "F9", "F10"}
				_ = t.hKBox[0].SetText(tmpList[0])
				_ = t.hKBox[1].SetText(tmpList[1])
				_ = t.hKBox[2].SetText(tmpList[2])
				_ = t.hKBox[3].SetText(tmpList[3])
			}},

			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: t.languageMap[language.OKStr], OnClicked: func() {
				M := make(map[string]bool)
				for _, v := range tmpList {
					if M[v] {
						walk.MsgBox(dlg, t.languageMap[language.ErrWindowTitleStr], t.languageMap[language.SetHotKeyErrMessageStr], walk.MsgBoxIconInformation)
						return
					} else {
						M[v] = true
					}
				}
				t.hKList = tmpList
				dlg.Accept()
			}},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: t.languageMap[language.CancelStr], OnClicked: func() { dlg.Cancel() }},
		},
	}.Run(t.mw)

	return cmd, err
}

// --------------------------------------- 订阅事件 ----------------------------------------------

// 订阅错误事件
func (t *SystemT) subShowError(data interface{}) (err error) {
	d := data.(events.ServerErrorData)

	if t.errorEdit == nil {
		return
	}

	return t.errorEdit.SetText(d.ErrInfo)
}

// 订阅状态变动事件
func (t *SystemT) subServerChange(data interface{}) (err error) {
	d := data.(events.ServerChangeData)

	for t.statusEdit == nil {
		time.Sleep(10 * time.Millisecond)
	}

	if err = t.statusEdit.SetText(string(d.Status)); err != nil {
		return
	}

	return

}
