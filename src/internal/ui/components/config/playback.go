package uiConfig

import (
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/event_topic"
	"math"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

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
	eventCenter.Event.Register(event_topic.LanguageChange, c.languageHandler)
	eventCenter.Event.Register(event_topic.ConfigChange, c.playbackRemainTimesHandler)
}

// 语言变动回调
func (c *playbackConfig) languageHandler(data interface{}) (err error) {
	for c.playbackTimesLabel == nil || c.currentTimesLabel == nil || c.speedLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	tryPublishErr(c.playbackTimesLabel.SetText(language.PlayBackTimesLabelStr.ToString()))
	tryPublishErr(c.currentTimesLabel.SetText(language.CurrentTimesLabelStr.ToString()))
	tryPublishErr(c.speedLabel.SetText(language.SpeedLabelStr.ToString()))
	return
}

// 回放剩余次数变动回调
func (c *playbackConfig) playbackRemainTimesHandler(data interface{}) (err error) {
	if confData, ok := data.(event_topic.ConfigChangeData); ok {
		if confData.Key == conf.KeyPlaybackTimes.GetKey() {
			// 获取回放次数
			var value = conf.KeyPlaybackTimes.GetValue()
			if value, ok := value.(int64); ok {
				tryPublishErr(c.currentTimesEdit.SetValue(float64(value)))
			}
		}
	}
	return
}

// 设置回放速度
func (c *playbackConfig) setSpeed() {
	tryPublishErr(c.speedEdit.SetValue(math.Pow(2, float64(c.speedSli.Value()-5))))

	conf.KeyPlaybackSpeed.SetValue(c.speedEdit.Value())
}

// 设置回放次数
func (c *playbackConfig) setPlaybackTimes() {
	conf.KeyPlaybackTimes.SetValue(int64(c.playbackTimesEdit.Value()))
}
