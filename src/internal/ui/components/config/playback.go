package component_config

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	hk "KeyMouseSimulation/pkg/hotkey"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"fmt"
	"math"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
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

	widget []declarative.Widget
}

func (c *playbackConfig) init() {
	var (
		playBackTimes = conf.PlaybackTimesConf.GetValue()
		speed         = conf.PlaybackSpeedConf.GetValue()
		speedSliValue = int(math.Log2(speed) + 5)
	)
	c.widget = []declarative.Widget{
		//速度设置
		declarative.Label{AssignTo: &c.speedLabel, ColumnSpan: 3},
		declarative.NumberEdit{AssignTo: &c.speedEdit, ColumnSpan: 5, Value: speed, MinValue: 0.01, MaxValue: 100.00, Suffix: " X", Decimals: 4},
		declarative.Slider{AssignTo: &c.speedSli, ColumnSpan: 8, MinValue: 0, MaxValue: 10, Value: speedSliValue, ToolTipsHidden: true, OnValueChanged: c.setSpeed},

		//回放次数与当前次数
		declarative.Label{AssignTo: &c.playbackTimesLabel, ColumnSpan: 2},
		declarative.NumberEdit{AssignTo: &c.playbackTimesEdit, ColumnSpan: 3, MinValue: float64(math.MinInt64), MaxValue: float64(math.MaxInt64), Value: float64(playBackTimes),
			StretchFactor: 0, Suffix: " times", Decimals: 0, OnValueChanged: c.setPlaybackTimes},
		declarative.Label{AssignTo: &c.currentTimesLabel, ColumnSpan: 2},
		declarative.NumberEdit{AssignTo: &c.currentTimesEdit, StretchFactor: 0, ColumnSpan: 1},
	}

	// 注册回调
	c.register()

}
func (c *playbackConfig) disPlay() []declarative.Widget {
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

	// 补充热键展示
	c.speedLabelShow()
}
func (c *playbackConfig) speedLabelShow() {
	var show = hk.Center.ShowSign()
	var text = fmt.Sprintf("%s[ '%s' - '%s' ]", language.SpeedLabelStr.ToString(), show[enum.HotKeySpeedDown], show[enum.HotKeySpeedUp])
	uiComponent.TryPublishErr(c.speedLabel.SetText(text))
}

//  ---------------------------------- 注册回调 ----------------------------------

func (c *playbackConfig) register() {
	// 热键变动
	eventCenter.Event.Register(topic.HotKeySet, c.hotKeySetHandler)
	// 热键触发
	eventCenter.Event.Register(topic.HotKeyEffect, c.hotKeyEffectHandler)

	// 配置变动
	eventCenter.Event.Register(topic.ConfigChange, c.configChangeHandler)
}

// 热键变动
func (c *playbackConfig) hotKeySetHandler(data interface{}) (err error) {
	c.speedLabelShow()
	return
}

// 热键触发
func (c *playbackConfig) hotKeyEffectHandler(data interface{}) (err error) {
	var dataValue = data.(*topic.HotKeyEffectData)

	var speed = c.speedSli.Value()
	switch dataValue.HotKey {
	case enum.HotKeySpeedUp:
		speed += 1
	case enum.HotKeySpeedDown:
		speed -= 1
	}

	c.speedSli.SetValue(speed)
	uiComponent.TryPublishErr(c.speedEdit.SetValue(math.Pow(2, float64(speed-5))))
	conf.PlaybackSpeedConf.SetValue(c.speedEdit.Value())

	return
}

// 配置变动
func (c *playbackConfig) configChangeHandler(data interface{}) (err error) {
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
}
