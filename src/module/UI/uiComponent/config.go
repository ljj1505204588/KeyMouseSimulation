package uiComponent

import (
	gene "KeyMouseSimulation/common/GenTool"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/module/language"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"math"
	"sync"
	"time"
)

type ConfigManageT struct {
	mw **walk.MainWindow
	sync.Once

	configs []configI // 配置
}

func (t *ConfigManageT) Init() {
	t.configs = []configI{&fileConfig{}, &recordConfig{}, &playbackConfig{}}
	for _, conf := range t.configs {
		conf.init()
		conf.register()
	}
}

func (t *ConfigManageT) DisPlay(mw **walk.MainWindow) (res []Widget) {
	t.mw = mw
	t.Once.Do(t.Init)

	for _, conf := range t.configs {
		res = append(res, conf.disPlay()...)
	}
	return
}

// --------------------------------------- 文件 ----------------------------------------------

type configI interface {
	init()
	disPlay() []Widget
	register()
}
type fileConfig struct {
	fileNames     []string
	fileComponent component.FileControlI

	fileLabel *walk.Label
	fileBox   *walk.ComboBox

	widget []Widget
}

func (c *fileConfig) init() {
	c.widget = []Widget{
		Label{AssignTo: &c.fileLabel, ColumnSpan: 2},
		ComboBox{AssignTo: &c.fileBox, ColumnSpan: 6, Model: c.fileNames, Editable: true, OnCurrentIndexChanged: c.chooseFile},
	}

	c.fileComponent = component.FileControl
}
func (c *fileConfig) disPlay() []Widget {

	return c.widget
}

func (c *fileConfig) register() {
	language.Center.RegisterChange(c.languageHandler)
	c.fileComponent.FileChange(c.fileChangeHandler)
}

// 语言变动回调
func (c *fileConfig) languageHandler() {
	for c.fileLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}
	tryPublishErr(c.fileLabel.SetText(language.Center.Get(language.FileLabelStr)))
}

// chooseFile 选择文件
func (c *fileConfig) chooseFile() {
	tryPublishErr(c.fileComponent.Choose(c.fileBox.Text()))
}

// 文件变动
func (c *fileConfig) fileChangeHandler(names []string, newFile []string) {
	if c.fileLabel == nil || c.fileBox == nil {
		return
	}

	// 模式设置
	tryPublishErr(c.fileBox.SetModel(names))

	// 新文件设置
	if len(newFile) != 0 {
		tryPublishErr(c.fileComponent.Choose(newFile[0]))
		tryPublishErr(c.fileBox.SetText(newFile[0]))
	}

	// 文件删除重设
	if !gene.Contain(names, c.fileBox.Text()) {
		var current = ""
		if len(names) > 0 {
			current = names[0]
		}
		tryPublishErr(c.fileComponent.Choose(current))
		tryPublishErr(c.fileBox.SetText(current))
	}

	return
}

// --------------------------------------- 记录 ----------------------------------------------

type recordConfig struct {

	// 鼠标路径记录
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	widget []Widget
}

func (c *recordConfig) init() {
	c.widget = []Widget{
		//鼠标路径
		Label{AssignTo: &c.ifMouseTrackLabel, Text: language.Center.Get(language.MouseTrackStr), ColumnSpan: 4},
		CheckBox{AssignTo: &c.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged: c.setIfTrackMouseMoveClick},
	}
}
func (c *recordConfig) disPlay() []Widget {
	return c.widget
}
func (c *recordConfig) register() {
	language.Center.RegisterChange(c.languageHandler)
	component.RecordConfig.SetMouseTrackChange(false, c.ifTrackMouseMoveRegister)
}

// 语言变动回调
func (c *recordConfig) languageHandler() {
	for c.ifMouseTrackLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	tryPublishErr(c.ifMouseTrackLabel.SetText(language.Center.Get(language.MouseTrackStr)))
}

// 设置是否追踪鼠标移动路径
func (c *recordConfig) setIfTrackMouseMoveClick() {
	component.RecordConfig.SetMouseTrack(c.ifMouseTrackCheck.Checked())
}

// 鼠标移动路径变动回调
func (c *recordConfig) ifTrackMouseMoveRegister(track bool) {
	c.ifMouseTrackCheck.SetChecked(track)
}

// --------------------------------------- 回放 ----------------------------------------------

type playbackConfig struct {
	// 回放次数调整
	playbackTimesLabel *walk.Label
	playbackTimesEdit  *walk.NumberEdit
	currentTimesLabel  *walk.Label
	currentTimesEdit   *walk.NumberEdit

	// 速度
	speedLabel *walk.Label
	speedSli   *walk.Slider
	speedEdit  *walk.NumberEdit

	widget []Widget
}

func (c *playbackConfig) init() {
	c.widget = []Widget{
		//速度设置
		Label{AssignTo: &c.speedLabel, ColumnSpan: 2},
		NumberEdit{AssignTo: &c.speedEdit, ColumnSpan: 6, MinValue: 0.01, MaxValue: 100.00, Suffix: " X", Decimals: 4},
		Slider{AssignTo: &c.speedSli, ColumnSpan: 8, MinValue: 0, MaxValue: 10, Value: 5, ToolTipsHidden: true, OnValueChanged: c.setSpeed},

		//回放次数与当前次数
		Label{AssignTo: &c.playbackTimesLabel, ColumnSpan: 2},
		NumberEdit{AssignTo: &c.playbackTimesEdit, ColumnSpan: 3, MinValue: -1, MaxValue: float64(math.MaxInt64), Value: float64(1), StretchFactor: 0, Suffix: " times", Decimals: 0, OnValueChanged: c.setPlaybackTimes},
		Label{AssignTo: &c.currentTimesLabel, ColumnSpan: 2},
		NumberEdit{AssignTo: &c.currentTimesEdit, StretchFactor: 0, ColumnSpan: 1},
	}

}
func (c *playbackConfig) disPlay() []Widget {
	return c.widget
}

func (c *playbackConfig) register() {
	language.Center.RegisterChange(c.languageHandler)
	component.PlaybackConfig.SetPlaybackRemainTimesChange(false, c.playbackRemainTimesHandler)
}

// 语言变动回调
func (c *playbackConfig) languageHandler() {
	for c.playbackTimesLabel == nil || c.currentTimesLabel == nil || c.speedLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	tryPublishErr(c.playbackTimesLabel.SetText(language.Center.Get(language.PlayBackTimesLabelStr)))
	tryPublishErr(c.currentTimesLabel.SetText(language.Center.Get(language.CurrentTimesLabelStr)))
	tryPublishErr(c.speedLabel.SetText(language.Center.Get(language.SpeedLabelStr)))
}

// 设置回放速度
func (c *playbackConfig) setSpeed() {
	var err = c.speedEdit.SetValue(math.Pow(2, float64(c.speedSli.Value()-5)))
	tryPublishErr(err)

	component.PlaybackConfig.SetSpeed(c.speedEdit.Value())
}

// 设置回放次数
func (c *playbackConfig) setPlaybackTimes() {
	component.PlaybackConfig.SetPlaybackTimes(int64(c.playbackTimesEdit.Value()))
}

// 回放剩余次数变动回调
func (c *playbackConfig) playbackRemainTimesHandler(times int64) {
	var err = c.currentTimesEdit.SetValue(float64(times))
	tryPublishErr(err)
}
