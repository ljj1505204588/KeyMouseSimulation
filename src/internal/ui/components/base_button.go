package uiComponent

import (
	"KeyMouseSimulation/common/common"
	component "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/file"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/event_topic"
	"fmt"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// FunctionT 热键按钮
type FunctionT struct {
	sync.Once

	file file.FileControlI
	functionWailT
}

type functionWailT struct {
	mw      **walk.MainWindow
	buttons []hotKeyButton
	widget  []Widget
}

type hotKeyButton struct {
	name   enum.HotKey
	exec   func()
	hkExec func()

	*walk.PushButton
}

func (t *FunctionT) Init() {
	t.file = file.FileControl

	t.buttons = []hotKeyButton{
		{name: enum.HotKeyRecord, exec: t.recordButtonClick, hkExec: t.recordButtonClick},
		{name: enum.HotKeyPlayBack, exec: t.playbackButtonClick, hkExec: t.playbackButtonClick},
		{name: enum.HotKeyPause, exec: t.pauseButtonClick, hkExec: t.pauseButtonClick},
		{name: enum.HotKeyStop, exec: t.stopButtonClick, hkExec: t.simStopButtonClick},
	}

	var err error
	for i, but := range t.buttons {
		t.buttons[i].HotKeyI, err = component.NewHK(but.name, but.name.DefKey(), but.hkExec)
		common.MustNil(err)
		t.widget = append(t.widget, PushButton{AssignTo: &(t.buttons[i].PushButton), ColumnSpan: 4, OnClicked: but.exec})
	}

	eventCenter.Event.Register(event_topic.RecordFinish, t.recordFinishHandler)
}

func (t *FunctionT) DisPlay(mw **walk.MainWindow) []Widget {
	t.mw = mw
	t.Once.Do(t.Init)
	return t.widget
}

// --------------------------------------- 按钮功能 ----------------------------------------------

func (t *FunctionT) recordButtonClick() {
	t.publishButtonClick(enum.RecordButton, "")
}
func (t *FunctionT) playbackButtonClick() {
	t.publishButtonClick(enum.PlaybackButton, t.file.Current())
}
func (t *FunctionT) pauseButtonClick() {
	t.publishButtonClick(enum.PauseButton, "")
}
func (t *FunctionT) stopButtonClick() {
	t.publishButtonClick(enum.StopButton, "")
}
func (t *FunctionT) simStopButtonClick() {
	// 停止涉及弹窗，目前考虑是这样特殊实现
	for _, but := range t.buttons {
		if but.name == enum.HotKeyStop {
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
	var data = dataI.(event_topic.RecordFinishData)
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
	tryPublishErr(err)

	return fileName, cmd == walk.DlgCmdOK && err == nil
}

// LanguageChange 设置语言
func (t *FunctionT) LanguageChange(data interface{}) (err error) {

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	for _, but := range t.buttons {
		tryPublishErr(but.SetText(fmt.Sprintf("%s %s", component.Center.Get(but.name.Language), but.Key())))
	}

	reutrn
}

// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *FunctionT) publishButtonClick(button enum.Button, name string) {
	_ = eventCenter.Event.Publish(event_topic.ButtonClick, event_topic.ButtonClickData{
		Button: button,
		Name:   name,
	})
}
