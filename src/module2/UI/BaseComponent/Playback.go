package BaseComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module2/language"
	"KeyMouseSimulation/share/events"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"math"
	"time"
)

// PlaybackT 回放按钮
type PlaybackT struct {
	*BaseT

	//文件选择
	fileLabel *walk.Label

	//回放次数调整
	playbackTimesLabel *walk.Label
	playbackTimesEdit  *walk.NumberEdit
	currentTimesLabel  *walk.Label
	currentTimesEdit   *walk.NumberEdit

	//速度
	speedLabel *walk.Label
	speedSli   *walk.Slider
	speedEdit  *walk.NumberEdit

	widget []Widget
}

func (t *PlaybackT) Init(base *BaseT) {
	t.BaseT = base
	t.BaseT.registerChangeLanguage(t.changeLanguageHandler)

	t.widget = []Widget{

		//回放文件
		Label{AssignTo: &t.fileLabel, ColumnSpan: 2},
		ComboBox{AssignTo: &t.fileBox, ColumnSpan: 6, Model: t.fileNames, Editable: true}, // OnCurrentIndexChanged: t.setFileName},

		//速度设置
		Label{AssignTo: &t.speedLabel, ColumnSpan: 2},
		NumberEdit{AssignTo: &t.speedEdit, ColumnSpan: 6, MinValue: 0.01, MaxValue: 100.00, Suffix: " X", Decimals: 4},
		Slider{AssignTo: &t.speedSli, ColumnSpan: 8, MinValue: 0, MaxValue: 10, Value: 5, ToolTipsHidden: true, OnValueChanged: t.setSpeed},

		//回放次数与当前次数
		Label{AssignTo: &t.playbackTimesLabel, ColumnSpan: 2},
		NumberEdit{AssignTo: &t.playbackTimesEdit, ColumnSpan: 3, MinValue: -1, MaxValue: float64(math.MaxInt64), Value: float64(1), StretchFactor: 0, Suffix: " times", Decimals: 0, OnValueChanged: t.setPlaybackTimes},
		Label{AssignTo: &t.currentTimesLabel, ColumnSpan: 2},
		NumberEdit{AssignTo: &t.currentTimesEdit, StretchFactor: 0, ColumnSpan: 1},
	}

	eventCenter.Event.Register(events.ServerChange, t.subServerChange)
}

func (t *PlaybackT) DisPlay() []Widget {
	return t.widget
}

// --------------------------------------- 基础功能 ----------------------------------------------

func (t *PlaybackT) initCheck() bool {
	for _, per := range []*walk.NumberEdit{
		t.playbackTimesEdit,
		t.currentTimesEdit,
		t.speedEdit,
	} {
		if per == nil {
			return false
		}
	}

	for _, per := range []*walk.Label{
		t.fileLabel,
		t.playbackTimesLabel,
		t.currentTimesLabel,
		t.speedLabel,
	} {
		if per == nil {
			return false
		}
	}
	return t.speedSli != nil
}

// 设置回放速度
func (t *PlaybackT) setSpeed() {
	//	defer t.lockSelf()()
	//	_ = t.speedEdit.SetValue(math.Pow(2, float64(t.speedSli.Value()-5)))
	//	t.sc.SetSpeed(t.speedEdit.Value())
}

// 设置回放次数
func (t *PlaybackT) setPlaybackTimes() {
	//	defer t.lockSelf()()
	//	t.sc.SetPlaybackTimes(int(t.playbackTimesEdit.Value()))
}

// 设置语言
func (t *PlaybackT) changeLanguageHandler(typ language.LanguageTyp) {
	var m = language.LanguageMap[typ]

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	_ = t.fileLabel.SetText(m[language.FileLabelStr])
	_ = t.playbackTimesLabel.SetText(m[language.PlayBackTimesLabelStr])
	_ = t.currentTimesLabel.SetText(m[language.CurrentTimesLabelStr])
	_ = t.speedLabel.SetText(m[language.SpeedLabelStr])
}

// --------------------------------------- 订阅事件 ----------------------------------------------

// 订阅状态变动事件
func (t *PlaybackT) subServerChange(data interface{}) (err error) {
	//d := data.(events.ServerChangeData)
	//
	//for t.currentTimesEdit == nil || t.fileBox == nil {
	//	time.Sleep(10 * time.Millisecond)
	//}
	//
	//if err = t.currentTimesEdit.SetValue(float64(d.CurrentTimes)); err != nil {
	//	return
	//}
	//if d.FileNamesData.Change {
	//	t.fileNames = d.FileNamesData.FileNames
	//	if err = t.fileBox.SetModel(d.FileNamesData.FileNames); err != nil {
	//		return
	//	}
	//	if t.fileBox.Text() == "" && len(d.FileNamesData.FileNames) != 0 {
	//		t.sc.SetFileName(d.FileNamesData.FileNames[0])
	//		if err = t.fileBox.SetText(d.FileNamesData.FileNames[0]); err != nil {
	//			return
	//		}
	//	}
	//}
	return
}
