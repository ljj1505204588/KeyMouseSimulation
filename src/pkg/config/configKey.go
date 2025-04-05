package conf

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
)

// 预定义配置键
var (
	RecordHotKeyConf    = &configKeyT[string]{key: enum.RecordHotKeyConf, val: "F7"}
	PlayBackHotKeyConf  = &configKeyT[string]{key: enum.PlaybackHotKeyConf, val: "F8"}
	PauseHotKeyConf     = &configKeyT[string]{key: enum.PauseHotKeyConf, val: "F9"}
	StopHotKeyConf      = &configKeyT[string]{key: enum.StopHotKeyConf, val: "F10"}
	SpeedUpHotKeyConf   = &configKeyT[string]{key: enum.SpeedUpHotKeyConf, val: ">"}
	SpeedDownHotKeyConf = &configKeyT[string]{key: enum.SpeedDownHotKeyConf, val: "<"}

	RecordLen               = &configKeyT[int]{key: enum.RecordLenConf, val: 0}
	RecordMouseTrackConf    = &configKeyT[bool]{key: enum.RecordMouseTrackConf, val: true}
	PlaybackSpeedConf       = &configKeyT[float64]{key: enum.PlaybackSpeedConf, val: 1.0}
	PlaybackTimesConf       = &configKeyT[int64]{key: enum.PlaybackTimesConf, val: int64(1)}
	PlaybackRemainTimesConf = &configKeyT[int64]{key: enum.PlaybackRemainTimesConf, val: int64(1)}
	LanguageConf            = &configKeyT[string]{key: enum.LanguageConf, val: "zh_CN"}
)

// configKeyT 泛型配置键
type configKeyT[T any] struct {
	key enum.ConfEnum
	val T
}

// GetKey 获取键
func (c *configKeyT[T]) GetKey() enum.ConfEnum {
	return c.key
}

// GetValue 获取值
func (c *configKeyT[T]) GetValue() T {
	return c.val
}

// SetValue 设置值
func (c *configKeyT[T]) SetValue(value T) {
	c.val = value

	// 发布配置改变事件
	eventCenter.Event.ASyncPublish(topic.ConfigChange, &topic.ConfigChangeData{Key: c.key, Value: c.val})
}
