package component_config

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
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

	// 注册回调
	eventCenter.Event.Register(topic.ConfigChange, func(data interface{}) (err error) {
		var value = data.(*topic.ConfigChangeData)

		switch value.Key {
		case enum.PlaybackTimesConf:
			// 回放次数
			uiComponent.TryPublishErr(c.playbackTimesEdit.SetValue(float64(value.Value.(int64))))
		case enum.PlaybackRemainTimesConf:
			// 剩余回放次数
			uiComponent.TryPublishErr(c.currentTimesEdit.SetValue(float64(value.Value.(int64))))
		}

		return
	})

}
func (c *playbackConfig) disPlay() []Widget {
	return c.widget
}

// 设置回放速度
func (c *playbackConfig) setSpeed() {
	uiComponent.TryPublishErr(c.speedEdit.SetValue(math.Pow(2, float64(c.speedSli.Value()-5))))
	conf.PlaybackSpeedConf.SetValue(c.speedEdit.Value())
}

// 设置回放次数
func (c *playbackConfig) setPlaybackTimes() {
	conf.PlaybackTimesConf.SetValue(int64(c.playbackTimesEdit.Value()))
}

// 语言变动回调
func (c *playbackConfig) languageChange() {
	for c.playbackTimesLabel == nil || c.currentTimesLabel == nil || c.speedLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	uiComponent.TryPublishErr(c.playbackTimesLabel.SetText(language.PlayBackTimesLabelStr.ToString()))
	uiComponent.TryPublishErr(c.currentTimesLabel.SetText(language.CurrentTimesLabelStr.ToString()))
	uiComponent.TryPublishErr(c.speedLabel.SetText(language.SpeedLabelStr.ToString()))
}
