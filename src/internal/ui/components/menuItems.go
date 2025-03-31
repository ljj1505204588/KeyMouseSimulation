package uiComponent

import (
	"KeyMouseSimulation/common/gene"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	hk "KeyMouseSimulation/pkg/hotkey"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/event_topic"
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
				Action{Text: string(enum.English), OnTriggered: func() {
					_ = eventCenter.Event.Publish(event_topic.LanguageChange, &event_topic.LanguageChangeData{
						Typ: enum.English,
					})
				}},
				Action{Text: string(enum.Chinese), OnTriggered: func() {
					_ = eventCenter.Event.Publish(event_topic.LanguageChange, &event_topic.LanguageChangeData{
						Typ: enum.Chinese,
					})
				}},
			},
			},
		}},
		Menu{AssignActionTo: &t.helpMenu, Items: []MenuItem{
			Action{AssignTo: &t.aboutAction, OnTriggered: t.showAboutBoxAction},
		}},
	}

}

// 设置热键弹窗
func (t *MenuItemT) setHotKeyPop() {
	type hkSetT struct {
		name   enum.HotKey
		box    *walk.ComboBox
		widget []Widget
	}

	var resetMethod []func()

	var hkSets []hkSetT
	for _, key := range enum.TotalHotkey() {
		var set = hkSetT{name: key, box: &walk.ComboBox{}}
		set.widget = []Widget{
			Label{Text: hk.Show(key), ColumnSpan: 1},
			ComboBox{AssignTo: &set.box, ColumnSpan: 1, Model: keyMouTool.VKCodeStringKeys, Editable: true, Value: hk.ShowSign(key), OnCurrentIndexChanged: func() {
				setM[set.hk] = set.box.Text()
			}},
		}

		hkSets = append(hkSets, set)

		var defKey = key.DefaultKey()
		resetMethod = append(resetMethod, func() {
			_ = set.box.SetText(defKey)
			setM[key] = defKey
		})
		setM[key] = key.Key() // 校验重复性
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
		PushButton{ColumnSpan: 4, Text: language.ResetStr.ToString(), OnClicked: func() {
			for _, reset := range resetMethod {
				reset()
			}
		}},

		PushButton{AssignTo: &t.acceptPB, ColumnSpan: 2, Text: language.OKStr.ToString(), OnClicked: func() {
			var keys []string
			for _, key := range setM {
				keys = append(keys, key)
			}

			if len(gene.RemoveDuplicate(keys)) != len(keys) {
				walk.MsgBox(dlg, language.ErrWindowTitleStr.ToString(), language.SetHotKeyErrMessageStr.ToString(), walk.MsgBoxIconInformation)
				return
			}

			dlg.Accept()
		}},
		PushButton{AssignTo: &t.cancelPB, ColumnSpan: 2, Text: language.CancelStr.ToString(), OnClicked: func() { dlg.Cancel() }},
	}...)

	cmd, _ := Dialog{AssignTo: &dlg, Title: language.SetHotKeyWindowTitleStr.ToString(),
		DefaultButton: &t.acceptPB, CancelButton: &t.cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: hkWidget,
	}.Run(*t.mw)

	if cmd == walk.DlgCmdOK {
		//tryPublishErr(component.MulSetKey(setM))
		//component.Center.Refresh()
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

// LanguageChange 设置语言
func (t *MenuItemT) LanguageChange(data interface{}) (err error) {
	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	tryPublishErr(t.settingMenu.SetText(language.MenuSettingStr.ToString()))
	tryPublishErr(t.languageMenu.SetText(language.MenuItemLanguageStr.ToString()))
	tryPublishErr(t.setHotkeyAction.SetText(language.ActionSetHotKeyStr.ToString()))
	tryPublishErr(t.helpMenu.SetText(language.MenuHelpStr.ToString()))
	tryPublishErr(t.aboutAction.SetText(language.ActionAboutStr.ToString()))

	return
}

// 系统信息
func (t *MenuItemT) showAboutBoxAction() {
	walk.MsgBox(*t.mw, language.AboutWindowTitleStr.ToString(), language.AboutMessageStr.ToString(), walk.MsgBoxIconInformation)
}
