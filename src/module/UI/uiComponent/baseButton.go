package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/commonTool"
	enum2 "KeyMouseSimulation/common/share/enum"
	"KeyMouseSimulation/common/share/events"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/module/language"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"sync"
	"time"
)

// FunctionT 热键按钮
type FunctionT struct {
	sync.Once

	file component.FileControlI
	functionWailT
}

type functionWailT struct {
	mw      **walk.MainWindow
	buttons []hotKeyButton
	widget  []Widget
}

type hotKeyButton struct {
	name   enum2.HotKey
	exec   func()
	hkExec func()

	*walk.PushButton
	component.HotKeyI
}

func (t *FunctionT) Init() {
	t.file = component.FileControl

	t.buttons = []hotKeyButton{
		{name: enum2.HotKeyRecord, exec: t.recordButtonClick, hkExec: t.recordButtonClick},
		{name: enum2.HotKeyPlayBack, exec: t.playbackButtonClick, hkExec: t.playbackButtonClick},
		{name: enum2.HotKeyPause, exec: t.pauseButtonClick, hkExec: t.pauseButtonClick},
		{name: enum2.HotKeyStop, exec: t.stopButtonClick, hkExec: t.simStopButtonClick},
	}

	var err error
	for i, but := range t.buttons {
		t.buttons[i].HotKeyI, err = component.NewHK(but.name, but.name.DefKey(), but.hkExec)
		commonTool.MustNil(err)
		t.widget = append(t.widget, PushButton{AssignTo: &(t.buttons[i].PushButton), ColumnSpan: 4, OnClicked: but.exec})
	}

	language.Center.RegisterChange(t.changeLanguageHandler)
	eventCenter.Event.Register(events.RecordFinish, t.recordFinishHandler)
}

func (t *FunctionT) DisPlay(mw **walk.MainWindow) []Widget {
	t.mw = mw
	t.Once.Do(t.Init)
	return t.widget
}

// --------------------------------------- 按钮功能 ----------------------------------------------

func (t *FunctionT) recordButtonClick() {
	t.publishButtonClick(enum2.RecordButton, "")
}
func (t *FunctionT) playbackButtonClick() {
	t.publishButtonClick(enum2.PlaybackButton, t.file.Current())
}
func (t *FunctionT) pauseButtonClick() {
	t.publishButtonClick(enum2.PauseButton, "")
}
func (t *FunctionT) stopButtonClick() {
	t.publishButtonClick(enum2.StopButton, "")
}
func (t *FunctionT) simStopButtonClick() {
	// 停止涉及弹窗，目前考虑是这样特殊实现
	for _, but := range t.buttons {
		if but.name == enum2.HotKeyStop {
			(*t.mw).WndProc(but.Handle(), win.WM_LBUTTONDOWN, 0, 0)
			(*t.mw).WndProc(but.Handle(), win.WM_LBUTTONUP, 0, 0)
		}
	}
}

// --------------------------------------- 基础功能 ----------------------------------------------

func (t *FunctionT) initCheck() bool {
	for _, per := range t.buttons {
		if per.PushButton == nil {
			return false
		}
	}

	return true
}

func (t *FunctionT) recordFinishHandler(dataI any) error {
	var data = dataI.(events.RecordFinishData)
	if fileName, ok := t.setFileName(); ok {
		t.file.Save(fileName, data.Notes)
	}
	return nil
}

// 设置文件名称
func (t *FunctionT) setFileName() (fileName string, ok bool) {
	var nameEdit *walk.LineEdit
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	cmd, err := Dialog{AssignTo: &dlg, Title: language.Center.Get(language.SetFileWindowTitleStr),
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			TextLabel{Text: language.Center.Get(language.SetFileLabelStr), ColumnSpan: 4},
			LineEdit{AssignTo: &nameEdit, ColumnSpan: 4, OnTextChanged: func() { fileName = nameEdit.Text() }},
			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.Center.Get(language.OKStr), OnClicked: func() { dlg.Accept() }},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.Center.Get(language.CancelStr), OnClicked: func() { dlg.Cancel() }},
		},
		Enabled: true,
	}.Run(*t.mw)
	tryPublishErr(err)

	return fileName, cmd == walk.DlgCmdOK && err == nil
}

// 设置语言
func (t *FunctionT) changeLanguageHandler() {

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	for _, but := range t.buttons {
		tryPublishErr(but.SetText(fmt.Sprintf("%s %s", language.Center.Get(but.name.Language()), but.Key())))
	}
}

// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *FunctionT) publishButtonClick(button enum2.Button, name string) {
	_ = eventCenter.Event.Publish(events.ButtonClick, events.ButtonClickData{
		Button: button,
		Name:   name,
	})
}
