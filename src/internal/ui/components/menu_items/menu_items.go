package component_menu_items

import (
	"KeyMouseSimulation/common/gene"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	"KeyMouseSimulation/internal/ui/components"
	eventCenter "KeyMouseSimulation/pkg/event"
	hk "KeyMouseSimulation/pkg/hotkey"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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
			t.hotKeyInit(),
			t.languageInit(),
		}},
		t.aboutInit(),
	}

}

// ---------------------------------------- 热键设置 ----------------------------------------

// 热键初始化
func (t *MenuItemT) hotKeyInit() MenuItem {
	return Action{AssignTo: &t.setHotkeyAction, OnTriggered: t.setHotKeyPop}
}

// 设置热键弹窗
func (t *MenuItemT) setHotKeyPop() {
	var (
		name        = hk.Center.Show()
		defaultSign = hk.Center.DefaultShowSign()
		sign        = hk.Center.ShowSign()
	)

	var mulBox mulHkBoxT
	for _, key := range enum.TotalHotkey() {
		var hkBox = &hkBoxT{
			name: key,
			box:  &walk.ComboBox{},
		}
		hkBox.widget = []Widget{
			Label{Text: name[key], ColumnSpan: 1},
			ComboBox{AssignTo: &hkBox.box, ColumnSpan: 1, Model: keyMouTool.VKCodeStringKeys, Editable: true, Value: sign[key], OnCurrentIndexChanged: func() {
				hkBox.setSign = hkBox.box.Text()
			}},
		}
		mulBox = append(mulBox, hkBox)
	}

	// 拼接页面
	var hkWidget = mulBox.getWidgets()

	var dlg *walk.Dialog
	hkWidget = append(hkWidget, []Widget{
		// 重置按钮
		PushButton{ColumnSpan: 4, Text: language.ResetStr.ToString(), OnClicked: func() {
			mulBox.resetSign(defaultSign)
		}},

		// 确认按钮
		PushButton{AssignTo: &t.acceptPB, ColumnSpan: 2, Text: language.OKStr.ToString(), OnClicked: func() {
			var keys = mulBox.getSigns()

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
		uiComponent.TryPublishErr(eventCenter.Event.Publish(topic.HotKeySet, &topic.HotKeySetData{
			Set: mulBox.getHotKeySet(),
		}))
	}
}

type hkBoxT struct {
	name    enum.HotKey
	box     *walk.ComboBox
	setSign string
	widget  []Widget
}

type mulHkBoxT []*hkBoxT

func (m mulHkBoxT) getWidgets() (res []Widget) {
	for _, box := range m {
		res = append(res, box.widget...)
	}
	return
}
func (m mulHkBoxT) resetSign(defSign map[enum.HotKey]string) {
	for _, hkBox := range m {
		hkBox.setSign = defSign[hkBox.name]
		_ = hkBox.box.SetText(hkBox.setSign)
	}
}
func (m mulHkBoxT) getSigns() (keys []string) {
	for _, box := range m {
		keys = append(keys, box.setSign)
	}
	return
}
func (m mulHkBoxT) getHotKeySet() (set map[enum.HotKey]keyMouTool.VKCode) {
	set = make(map[enum.HotKey]keyMouTool.VKCode)
	for _, box := range m {
		set[box.name] = keyMouTool.VKCodeStringMap[box.setSign]
	}
	return
}

// ---------------------------------------- 语言设置 ----------------------------------------

// 语言初始化
func (t *MenuItemT) languageInit() MenuItem {
	return Menu{AssignActionTo: &t.languageMenu, Items: []MenuItem{
		Action{Text: string(enum.English), OnTriggered: func() {
			_ = eventCenter.Event.Publish(topic.LanguageChange, &topic.LanguageChangeData{
				Typ: enum.English,
			})
		}},
		Action{Text: string(enum.Chinese), OnTriggered: func() {
			_ = eventCenter.Event.Publish(topic.LanguageChange, &topic.LanguageChangeData{
				Typ: enum.Chinese,
			})
		}},
	},
	}
}

// LanguageChange 设置语言
func (t *MenuItemT) LanguageChange(data interface{}) (err error) {
	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	uiComponent.TryPublishErr(t.settingMenu.SetText(language.MenuSettingStr.ToString()))         // 设置
	uiComponent.TryPublishErr(t.languageMenu.SetText(language.MenuItemLanguageStr.ToString()))   // 语言设置
	uiComponent.TryPublishErr(t.setHotkeyAction.SetText(language.ActionSetHotKeyStr.ToString())) // 热键设置
	uiComponent.TryPublishErr(t.helpMenu.SetText(language.MenuHelpStr.ToString()))               // 帮助
	uiComponent.TryPublishErr(t.aboutAction.SetText(language.ActionAboutStr.ToString()))         // 关于

	return
}

// 初始化校验
func (t *MenuItemT) initCheck() bool {
	for _, per := range []*walk.Action{
		t.settingMenu, t.setHotkeyAction, t.languageMenu,
		t.helpMenu, t.aboutAction,
	} {
		if per == nil {
			return false
		}
	}
	return true
}

// ---------------------------------------- 关于信息 ----------------------------------------

// 关于初始化
func (t *MenuItemT) aboutInit() MenuItem {
	return Menu{AssignActionTo: &t.helpMenu, Items: []MenuItem{
		Action{AssignTo: &t.aboutAction, OnTriggered: t.showAboutBoxAction},
	}}
}

// 系统信息
func (t *MenuItemT) showAboutBoxAction() {
	walk.MsgBox(*t.mw, language.AboutWindowTitleStr.ToString(), language.AboutMessageStr.ToString(), walk.MsgBoxIconInformation)
}
