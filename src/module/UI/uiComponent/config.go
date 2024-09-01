package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/share/events"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"math"
	"sync"
	"time"
)

type ConfigT struct {
	mw *walk.MainWindow
	sync.Once

	// 文件选择
	basePath      string
	fileNames     []string
	fileComponent component.FileControlI

	configWalk
	widget []Widget
}

type configWalk struct {
	// 文件选择
	fileLabel *walk.Label
	fileBox   *walk.ComboBox

	// 鼠标路径记录
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	// 回放次数调整
	playbackTimesLabel *walk.Label
	playbackTimesEdit  *walk.NumberEdit
	currentTimesLabel  *walk.Label
	currentTimesEdit   *walk.NumberEdit

	// 速度
	speedLabel *walk.Label
	speedSli   *walk.Slider
	speedEdit  *walk.NumberEdit
}

func (t *ConfigT) Init() {
	language.Center.RegisterChange(t.changeLanguageHandler)

	t.fileComponent = component.FileControl
	t.widget = []Widget{
		//鼠标路径
		Label{AssignTo: &t.ifMouseTrackLabel, ColumnSpan: 4},
		CheckBox{AssignTo: &t.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged: t.setIfTrackMouseMoveClick},

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

	eventCenter.Event.Register(events.ServerConfigChange, t.subServerChange)
	t.fileComponent.FileChange(t.fileChangeHandler)
}

func (t *ConfigT) DisPlay(mw *walk.MainWindow) []Widget {
	t.mw = mw
	t.Once.Do(t.Init)

	return t.widget
}

// 文件变动
func (t *ConfigT) fileChangeHandler(names []string, newFile []string) {
	if !t.initCheck() {
		return
	}

	if err := t.fileBox.SetModel(names); err != nil {
		tryPublishErr(err)
		return
	}

	if len(newFile) != 0 {
		if err := t.fileBox.SetText(newFile[0]); err != nil {
			tryPublishErr(err)
			return
		}
	}

	return
}

// 设置是否追踪鼠标移动路径
func (t *ConfigT) setIfTrackMouseMoveClick() {
	// defer t.lockSelf()()
	// t.sc.SetIfTrackMouseMove(t.ifMouseTrackCheck.Checked())
}

// --------------------------------------- 基础功能 ----------------------------------------------

func (t *ConfigT) initCheck() bool {
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
func (t *ConfigT) setSpeed() {
	//	defer t.lockSelf()()
	//	_ = t.speedEdit.SetValue(math.Pow(2, float64(t.speedSli.Value()-5)))
	//	t.sc.SetSpeed(t.speedEdit.Value())
}

// 设置回放次数
func (t *ConfigT) setPlaybackTimes() {
	//	defer t.lockSelf()()
	//	t.sc.SetPlaybackTimes(int(t.playbackTimesEdit.Value()))
}

// 设置语言
func (t *ConfigT) changeLanguageHandler() {

	for !t.initCheck() {
		time.Sleep(10 * time.Millisecond)
	}

	_ = t.fileLabel.SetText(language.Center.Get(language.FileLabelStr))
	_ = t.playbackTimesLabel.SetText(language.Center.Get(language.PlayBackTimesLabelStr))
	_ = t.currentTimesLabel.SetText(language.Center.Get(language.CurrentTimesLabelStr))
	_ = t.speedLabel.SetText(language.Center.Get(language.SpeedLabelStr))
}

// --------------------------------------- 订阅事件 ----------------------------------------------

// 订阅状态变动事件
func (t *ConfigT) subServerChange(data interface{}) (err error) {
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
