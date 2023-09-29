package BaseComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"time"
)

// FunctionT 热键按钮
type FunctionT struct {
	*BaseT

	recordButton      *walk.PushButton
	playbackButton    *walk.PushButton
	pauseButton       *walk.PushButton
	stopButton        *walk.PushButton
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	hotKeyHandler map[string]func()
	widget        []Widget
}

func (t *FunctionT) Init(base *BaseT) {
	t.BaseT = base

	t.widget = []Widget{
		PushButton{AssignTo: &t.recordButton, ColumnSpan: 4, OnClicked: t.recordButtonClick},
		PushButton{AssignTo: &t.playbackButton, ColumnSpan: 4, OnClicked: t.playbackButtonClick},
		PushButton{AssignTo: &t.pauseButton, ColumnSpan: 4, OnClicked: t.pauseButtonClick},
		PushButton{AssignTo: &t.stopButton, ColumnSpan: 4, OnClicked: t.stopButtonClick},
		//鼠标路径
		Label{AssignTo: &t.ifMouseTrackLabel, ColumnSpan: 4},
		CheckBox{AssignTo: &t.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged: t.setIfTrackMouseMoveClick},
	}

	t.Register()
}

func (t *FunctionT) DisPlay() []Widget {
	return t.widget
}

// --------------------------------------- 按钮功能 ----------------------------------------------

func (t *FunctionT) recordButtonClick() {
	defer t.lockSelf()()

	t.publishButtonClick(enum.RecordButton, "")
}
func (t *FunctionT) playbackButtonClick() {
	defer t.lockSelf()()

	t.publishButtonClick(enum.PlaybackButton, t.BaseT.fileBox.Text())
}
func (t *FunctionT) pauseButtonClick() {
	defer t.lockSelf()()

	t.publishButtonClick(enum.PauseButton, "")
}
func (t *FunctionT) stopButtonClick() {
	defer t.lockSelf()()

	t.publishButtonClick(enum.StopButton, "")
}

// --------------------------------------- 基础功能 ----------------------------------------------

func (t *FunctionT) initCheck() bool {
	for _, per := range []*walk.PushButton{
		t.recordButton,
		t.playbackButton,
		t.pauseButton,
		t.stopButton,
	} {
		if per == nil {
			return false
		}
	}

	return t.ifMouseTrackLabel != nil && t.ifMouseTrackCheck != nil
}

// 设置文件名称
func (t *FunctionT) setFileName() {
	var nameEdit *walk.LineEdit
	filename := ""
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var m = t.languageMap

	cmd, err := Dialog{AssignTo: &dlg, Title: m[language.SetFileWindowTitleStr],
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			TextLabel{Text: m[language.SetFileLabelStr], ColumnSpan: 4},
			LineEdit{AssignTo: &nameEdit, ColumnSpan: 4, OnTextChanged: func() { filename = nameEdit.Text() }},
			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: m[language.OKStr], OnClicked: func() { dlg.Accept() }},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: m[language.CancelStr], OnClicked: func() { dlg.Cancel() }},
		},
		Enabled: true,
	}.Run(t.mw)

	if cmd == walk.DlgCmdOK && err == nil {
		t.publishButtonClick(enum.SaveFileButton, filename)
	}
}

// 设置是否追踪鼠标移动路径
func (t *FunctionT) setIfTrackMouseMoveClick() {
	// defer t.lockSelf()()
	// t.sc.SetIfTrackMouseMove(t.ifMouseTrackCheck.Checked())
}

// 设置语言
func (t *FunctionT) changeLanguageHandler(typ language.LanguageTyp) {
	var m = language.LanguageMap[typ]

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	_ = t.recordButton.SetText(m[language.RecordStr] + " " + t.hKList[0])
	_ = t.playbackButton.SetText(m[language.PlaybackStr] + " " + t.hKList[1])
	_ = t.pauseButton.SetText(m[language.PauseStr] + " " + t.hKList[2])
	_ = t.stopButton.SetText(m[language.StopStr] + " " + t.hKList[3])
	_ = t.ifMouseTrackLabel.SetText(m[language.MouseTrackStr])
}

// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *FunctionT) publishButtonClick(button enum.Button, name string) {
	_ = eventCenter.Event.Publish(events.ButtonClick, events.ButtonClickData{
		Button: button,
		Name:   name,
	})

}

func (t *FunctionT) Register() {
	t.BaseT.registerChangeLanguage(t.changeLanguageHandler)

	// 停止涉及弹窗，目前考虑是这样特殊实现
	var simButtonStopClick = func() {
		t.mw.WndProc(t.stopButton.Handle(), win.WM_LBUTTONDOWN, 0, 0)
		t.mw.WndProc(t.stopButton.Handle(), win.WM_LBUTTONUP, 0, 0)
	}

	t.hotKeyHandler = map[string]func(){
		t.BaseT.hKList[0]: t.recordButtonClick,
		t.BaseT.hKList[1]: t.playbackButtonClick,
		t.BaseT.hKList[2]: t.pauseButtonClick,
		t.BaseT.hKList[3]: simButtonStopClick,
	}

	eventCenter.Event.Register(events.ServerHotKeyDown, t.hotKeyDownHandler)
	eventCenter.Event.Register(events.FileScanNewFile, t.fileChangeHandler)
	eventCenter.Event.Register(events.RecordFinish, t.recordFinishHandler)
}

// 记录结束
func (t *FunctionT) recordFinishHandler(data interface{}) (err error) {
	//d := data.(events.RecordFinishData)

	t.setFileName()
	return
}

// 热键按下
func (t *FunctionT) hotKeyDownHandler(data interface{}) (err error) {
	d := data.(events.ServerHotKeyDownData)

	if f, ok := t.hotKeyHandler[d.Key]; ok {
		f()
	}

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
