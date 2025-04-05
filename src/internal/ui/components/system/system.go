package component_system

import (
	"KeyMouseSimulation/internal/server"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/topic"
	"strings"
	"sync"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
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

	widgets []declarative.Widget
}

func (t *SystemT) Init() {

	t.widgets = []declarative.Widget{
		//当前状态
		declarative.Label{AssignTo: &t.statusLabel, ColumnSpan: 2},
		declarative.LineEdit{AssignTo: &t.statusEdit, ColumnSpan: 6, ReadOnly: true},

		//错误信息
		declarative.Label{AssignTo: &t.errorLabel, ColumnSpan: 2},
		declarative.TextEdit{AssignTo: &t.errorEdit, ColumnSpan: 6, ReadOnly: true},
	}

	t.register()
}

func (t *SystemT) DisPlay(mw **walk.MainWindow) []declarative.Widget {
	t.mw = mw
	t.Once.Do(t.Init)

	return t.widgets
}

// --------------------------------------- 基础功能 ----------------------------------------------

// 初始化校验
func (t *SystemT) initCheck() bool {
	return t.statusLabel != nil && t.statusEdit != nil && t.errorLabel != nil && t.errorEdit != nil
}

// LanguageChange 设置语言
func (t *SystemT) LanguageChange(data interface{}) (err error) {
	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	_ = t.statusLabel.SetText(language.StatusLabelStr.ToString())
	_ = t.errorLabel.SetText(language.ErrorLabelStr.ToString())
	return
}

// --------------------------------------- 订阅事件 ----------------------------------------------

func (t *SystemT) register() {

	eventCenter.Event.Register(topic.ServerError, t.subShowError)
	eventCenter.Event.Register(topic.ServerStatus, t.subServerStatusChange)
}

// 订阅错误事件
func (t *SystemT) subShowError(dataI interface{}) (err error) {
	if t.errorEdit == nil {
		return
	}

	// 日志记录
	var data = dataI.(*topic.ServerErrorData)
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
	d := data.(*topic.ServerStatusChangeData)

	for t.statusEdit == nil {
		time.Sleep(10 * time.Millisecond)
	}

	var show = server.Svc.StatusShow(d.Status)
	if err = t.statusEdit.SetText(show); err != nil {
		return
	}

	return
}
