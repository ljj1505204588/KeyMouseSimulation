package component_base_button

import (
	"KeyMouseSimulation/internal/server"
	"KeyMouseSimulation/internal/ui/components"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	hk "KeyMouseSimulation/pkg/hotkey"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"fmt"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// FunctionT 热键按钮
type FunctionT struct {
	mw      **walk.MainWindow
	buttons []hotKeyButton
	widget  []Widget

	isSaving bool // 是否存储中
}

type hotKeyButton struct {
	name enum.HotKey
	exec func()

	*walk.PushButton
}

func (t *FunctionT) Init() {

	t.buttons = []hotKeyButton{
		{name: enum.HotKeyRecord, exec: t.buttonWithCheck(server.Svc.Record)},
		{name: enum.HotKeyPlayBack, exec: t.buttonWithCheck(server.Svc.PlayBack)},
		{name: enum.HotKeyPause, exec: t.buttonWithCheck(server.Svc.Pause)},
		{name: enum.HotKeyStop, exec: t.buttonWithCheck(t.stopButtonClick)},
	}

	for i, but := range t.buttons {
		t.widget = append(t.widget, PushButton{AssignTo: &(t.buttons[i].PushButton), ColumnSpan: 4, OnClicked: but.exec})
	}

	t.hotKeyRegister()
}

func (t *FunctionT) DisPlay(mw **walk.MainWindow) []Widget {
	t.mw = mw
	t.Init()
	return t.widget
}

// --------------------------------------- 按钮功能 ----------------------------------------------

// 按下前校验
func (t *FunctionT) buttonWithCheck(f func()) (check func()) {
	check = func() {
		if !t.isSaving {
			f()
		}
	}
	return check
}

// 暂停
func (t *FunctionT) stopButtonClick() {
	t.isSaving = true
	defer func() { t.isSaving = false }()

	// 记录停止，要存储文件
	if server.Svc.Stop() && conf.RecordLen.GetValue() != 0 {
		if fileName, ok := t.setFileName(); ok {
			server.Svc.Save(fileName)
		}
	}

	conf.RecordLen.SetValue(0) // todo 考虑这个是否合理
}

// 按钮显示调整
func (t *FunctionT) buttonShow() {
	var (
		show     = hk.Center.Show()
		showSign = hk.Center.ShowSign()
	)
	for _, but := range t.buttons {
		var key = but.name
		uiComponent.TryPublishErr(but.SetText(fmt.Sprintf("%s %s", show[key], showSign[key])))
	}
}

// --------------------------------------- 热键 ----------------------------------------------

func (t *FunctionT) hotKeyRegister() {

	// 热键变动
	eventCenter.Event.Register(topic.HotKeySet, func(data interface{}) (err error) {
		t.buttonShow()
		return
	})

	// 注册热键触发
	eventCenter.Event.Register(topic.HotKeyEffect, func(data interface{}) (err error) {

		var dataValue = data.(*topic.HotKeyEffectData)
		for _, buttons := range t.buttons {
			// 执行热键
			if buttons.name == dataValue.HotKey {
				buttons.exec()
			}
		}

		return
	})
}

// --------------------------------------- 基础功能 ----------------------------------------------

// 设置文件名称
func (t *FunctionT) setFileName() (fileName string, ok bool) {
	var nameEdit *walk.LineEdit
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	cmd, err := Dialog{AssignTo: &dlg, Title: language.SetFileWindowTitleStr.ToString(),
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			TextLabel{Text: language.SetFileLabelStr.ToString(), ColumnSpan: 4},
			LineEdit{AssignTo: &nameEdit, ColumnSpan: 4, OnTextChanged: func() { fileName = nameEdit.Text() }},
			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.OKStr.ToString(), OnClicked: func() { dlg.Accept() }},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.CancelStr.ToString(), OnClicked: func() { dlg.Cancel() }},
		},
		Enabled: true,
	}.Run(*t.mw)
	uiComponent.TryPublishErr(err)

	return fileName, cmd == walk.DlgCmdOK && err == nil
}

// LanguageChange 设置语言
func (t *FunctionT) LanguageChange(data interface{}) (err error) {

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	t.buttonShow()

	return
}

func (t *FunctionT) initCheck() bool {
	for _, per := range t.buttons {
		if per.PushButton == nil {
			return false
		}
	}

	return true
}
