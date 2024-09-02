package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/share/events"
	"KeyMouseSimulation/module/language"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
	"sync"
	"time"
)

// SystemT 系统按钮
type SystemT struct {
	mw **walk.MainWindow
	sync.Once

	//系统状态
	statusLabel *walk.Label
	statusEdit  *walk.LineEdit

	errorLabel *walk.Label
	errorEdit  *walk.TextEdit
	historyErr []string

	widgets []Widget
}

func (t *SystemT) Init() {

	t.widgets = []Widget{
		//当前状态
		Label{AssignTo: &t.statusLabel, ColumnSpan: 2},
		LineEdit{AssignTo: &t.statusEdit, ColumnSpan: 6, ReadOnly: true},

		//错误信息
		Label{AssignTo: &t.errorLabel, ColumnSpan: 2},
		TextEdit{AssignTo: &t.errorEdit, ColumnSpan: 6, ReadOnly: true},
	}

	t.register()
}

func (t *SystemT) DisPlay(mw **walk.MainWindow) []Widget {
	t.mw = mw
	t.Once.Do(t.Init)

	return t.widgets
}

// --------------------------------------- 基础功能 ----------------------------------------------

// 初始化校验
func (t *SystemT) initCheck() bool {
	return t.statusLabel != nil && t.statusEdit != nil && t.errorLabel != nil && t.errorEdit != nil
}

// 修改语言
func (t *SystemT) changeLanguageHandler() {
	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	_ = t.statusLabel.SetText(language.Center.Get(language.StatusLabelStr))
	_ = t.errorLabel.SetText(language.Center.Get(language.ErrorLabelStr))
}

// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *SystemT) register() {
	language.Center.RegisterChange(t.changeLanguageHandler)

	eventCenter.Event.Register(events.ServerError, t.subShowError)
	eventCenter.Event.Register(events.ServerStatus, t.subServerStatusChange)
}

// 订阅错误事件
func (t *SystemT) subShowError(dataI interface{}) (err error) {
	if t.errorEdit == nil {
		return
	}

	// 日志记录
	var data = dataI.(events.ServerErrorData)
	t.historyErr = append(t.historyErr, data.ErrInfo+" \r\n")

	var hisLen = len(t.historyErr)
	var textBuild = strings.Builder{}
	for i := hisLen - 1; i >= 0 && i >= hisLen-10; i-- {
		textBuild.WriteString(t.historyErr[i])
	}

	if len(t.historyErr) > 1000 {
		t.historyErr = t.historyErr[hisLen-100:]
	}

	return t.errorEdit.SetText(textBuild.String())
}

// 订阅状态变动事件
func (t *SystemT) subServerStatusChange(data interface{}) (err error) {
	d := data.(events.ServerStatusChangeData)

	for t.statusEdit == nil {
		time.Sleep(10 * time.Millisecond)
	}

	if err = t.statusEdit.SetText(string(d.Status)); err != nil {
		return
	}

	return
}
