package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/commonTool"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"errors"
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

	functionWailT
}

type functionWailT struct {
	mw      *walk.MainWindow
	buttons []hotKeyButton
	widget  []Widget
}

type hotKeyButton struct {
	name      enum.HotKey
	key       string
	disPlayNo language.DisPlay
	exec      func()

	*walk.PushButton
	component.HotKeyI
}

func (t *FunctionT) Init() {
	t.buttons = []hotKeyButton{
		{name: enum.HotKeyRecord, key: "F7", disPlayNo: language.RecordStr, exec: t.recordButtonClick},
		{name: enum.HotKeyPlayBack, key: "F8", disPlayNo: language.PlaybackStr, exec: t.playbackButtonClick},
		{name: enum.HotKeyPause, key: "F9", disPlayNo: language.PauseStr, exec: t.pauseButtonClick},
		{name: enum.HotKeyStop, key: "F10", disPlayNo: language.StopStr, exec: t.stopButtonClick},
	}

	var err error
	for _, but := range t.buttons {
		but.HotKeyI, err = component.NewHK(but.name, but.key, but.exec)
		commonTool.MustNil(err)
		t.widget = append(t.widget, PushButton{AssignTo: &(but.PushButton), ColumnSpan: 4, OnClicked: but.exec})
	}

	language.Center.RegisterChange(t.changeLanguageHandler)
	eventCenter.Event.Register(events.RecordFinish, t.recordFinishHandler)

	t.widget = []Widget{}

}

func (t *FunctionT) DisPlay(mw *walk.MainWindow) []Widget {
	t.mw = mw
	t.Once.Do(t.Init)
	return t.widget
}

// --------------------------------------- 按钮功能 ----------------------------------------------

func (t *FunctionT) recordButtonClick() {
	// todo 考虑要不要锁
	t.publishButtonClick(enum.RecordButton, "")
}
func (t *FunctionT) playbackButtonClick() {
	t.publishButtonClick(enum.PlaybackButton, t.BaseT.fileBox.Text())
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
			t.mw.WndProc(but.Handle(), win.WM_LBUTTONDOWN, 0, 0)
			t.mw.WndProc(but.Handle(), win.WM_LBUTTONUP, 0, 0)
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

// 设置文件名称
func (t *FunctionT) setFileName() {
	var nameEdit *walk.LineEdit
	filename := ""
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	cmd, err := Dialog{AssignTo: &dlg, Title: language.Center.Get(language.SetFileWindowTitleStr),
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			TextLabel{Text: language.Center.Get(language.SetFileLabelStr), ColumnSpan: 4},
			LineEdit{AssignTo: &nameEdit, ColumnSpan: 4, OnTextChanged: func() { filename = nameEdit.Text() }},
			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.Center.Get(language.OKStr), OnClicked: func() { dlg.Accept() }},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.Center.Get(language.CancelStr), OnClicked: func() { dlg.Cancel() }},
		},
		Enabled: true,
	}.Run(t.mw)

	if cmd == walk.DlgCmdOK && err == nil {
		t.publishButtonClick(enum.SaveFileButton, filename)
	}
}

// 设置语言
func (t *FunctionT) changeLanguageHandler() {

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	for _, but := range t.buttons {
		_ = but.SetText(fmt.Sprintf("%s %s", language.Center.Get(but.disPlayNo), but.key))
	}
}

// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *FunctionT) publishButtonClick(button enum.Button, name string) {
	_ = eventCenter.Event.Publish(events.ButtonClick, events.ButtonClickData{
		Button: button,
		Name:   name,
	})

}

// 记录结束
func (t *FunctionT) recordFinishHandler(data interface{}) (err error) {
	t.setFileName()
	return
}

// 文件变动
func (t *FunctionT) fileChangeHandler(data interface{}) (err error) {
	if !t.initCheck() {
		return errors.New("wait init success")
	}

	var info = data.(events.FileScanNewFileData)
	if err = t.fileBox.SetModel(info.FileList); err != nil {
		return
	}

	if len(info.NewFile) != 0 {
		if err = t.fileBox.SetText(info.NewFile[0]); err != nil {
			return
		}
	}

	return
}