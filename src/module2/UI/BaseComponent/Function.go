package BaseComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module2/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"sort"
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

	widget []Widget
}

func (t *FunctionT)Init(base *BaseT){
	t.BaseT = base
	t.BaseT.registerChangeLanguage(t.changeLanguageHandler)

	t.widget = []Widget{
		PushButton{AssignTo: &t.recordButton, ColumnSpan: 4,  OnClicked: t.recordButtonClick},
		PushButton{AssignTo: &t.playbackButton, ColumnSpan: 4,  OnClicked: t.playbackButtonClick},
		PushButton{AssignTo: &t.pauseButton, ColumnSpan: 4,  OnClicked: t.pauseButtonClick},
		PushButton{AssignTo: &t.stopButton, ColumnSpan: 4,  OnClicked: t.stopButtonClick},
		//鼠标路径
		Label{AssignTo: &t.ifMouseTrackLabel, ColumnSpan: 4},
		CheckBox{AssignTo: &t.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged:t.setIfTrackMouseMoveClick },
	}
	eventCenter.Event.Register(events.ServerHotKeyDown, t.subHotKeyDown)
}

func (t *FunctionT)DisPlay()[]Widget{
	return t.widget
}

// --------------------------------------- 按钮功能 ----------------------------------------------

func(t *FunctionT)recordButtonClick(){
	defer  t.lockSelf()()
	t.sc.Record()
}
func(t *FunctionT)playbackButtonClick(){
	defer  t.lockSelf()()
	t.sc.Playback()
}
func(t *FunctionT)pauseButtonClick(){
	defer  t.lockSelf()()
	_ = t.sc.Pause()
}
func(t *FunctionT)stopButtonClick(){
	defer  t.lockSelf()()

	_ = t.sc.Pause()
	fileName, cmd, _ := t.setFileName()
	if cmd == walk.DlgCmdOK {
		for _, v := range t.fileNames {
			if v == fileName {
				fileName += "-" + time.Now().String()
			}
		}
		t.sc.SetFileName(fileName)
		_ = t.fileBox.SetText(fileName)

		t.fileNames = append(t.fileNames, fileName)
		sort.Strings(t.fileNames)
		_ = t.fileBox.SetModel(t.fileNames)
	}
	 t.sc.Stop(true)

}

// --------------------------------------- 基础功能 ----------------------------------------------

func(t *FunctionT) initCheck()bool{
	for _,per := range[]*walk.PushButton{
		t.recordButton,
		t.playbackButton,
		t.pauseButton,
		t.stopButton,
	}{
		if per == nil {
			return false
		}
	}

	return t.ifMouseTrackLabel != nil && t.ifMouseTrackCheck != nil
}

// 设置文件名称
func (t *FunctionT) setFileName() (string, int, error) {
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
	}.Run(t.mw)

	return filename, cmd, err
}
// 设置是否追踪鼠标移动路径
func (t *FunctionT) setIfTrackMouseMoveClick() {
	defer  t.lockSelf()()
	t.sc.SetIfTrackMouseMove(t.ifMouseTrackCheck.Checked())
}
// 设置语言
func (t *FunctionT)changeLanguageHandler(typ language.LanguageTyp){
	var m = language.LanguageMap[typ]

	for !t.initCheck(){
		time.Sleep(10*time.Millisecond)
	}

	_ = t.recordButton.SetText(m[language.RecordStr]   + " " + t.hKList[0])
	_ = t.playbackButton.SetText(m[language.PlaybackStr]   + " " + t.hKList[1])
	_ = t.pauseButton.SetText(m[language.PauseStr]   + " " + t.hKList[2])
	_ = t.stopButton.SetText(m[language.StopStr]   + " " + t.hKList[3])
	_ = t.ifMouseTrackLabel.SetText(m[language.MouseTrackStr])
}


// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *FunctionT) subHotKeyDown(data interface{}) (err error) {
	d := data.(events.ServerHotKeyDownData)

	switch d.HotKey {
	case enum.HOT_KEY_PLAYBACK_START:
		 t.clickButton(t.playbackButton)
	case enum.HOT_KEY_RECORD_START:
		 t.clickButton(t.recordButton)
	case enum.HOT_KEY_PAUSE:
		 t.clickButton(t.pauseButton)
	case enum.HOT_KEY_STOP:
		 t.clickButton(t.stopButton)
	}
	return
}

func (t *FunctionT) clickButton(button *walk.PushButton) {
	defer t.lockSelf()()

	t.mw.WndProc(button.Handle(), win.WM_LBUTTONDOWN, 0, 0)
	t.mw.WndProc(button.Handle(), win.WM_LBUTTONUP, 0, 0)
}