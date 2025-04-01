package conf

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"fmt"
)

// 预定义配置键
var (
	RecordHotKeyConf     = &configKeyT[string]{key: enum.RecordHotKeyConf, val: "F9"}
	PlaybackHotKeyConf   = &configKeyT[string]{key: enum.PlaybackHotKeyConf, val: "F10"}
	PauseHotKeyConf      = &configKeyT[string]{key: enum.PauseHotKeyConf, val: "F11"}
	StopHotKeyConf       = &configKeyT[string]{key: enum.StopHotKeyConf, val: "F12"}
	RecordMouseTrackConf = &configKeyT[bool]{key: enum.RecordMouseTrackConf, val: true}
	PlaybackSpeedConf    = &configKeyT[float64]{key: enum.PlaybackSpeedConf, val: 1.0}
	PlaybackTimesConf    = &configKeyT[int64]{key: enum.PlaybackTimesConf, val: int64(1)}
	LanguageConf         = &configKeyT[string]{key: enum.LanguageConf, val: "zh_CN"}
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
func (c *configKeyT[T]) GetValue() any {
	return c.val
}

// SetValue 设置值
func (c *configKeyT[T]) SetValue(value any) error {
	val, ok := value.(T)
	if !ok {
		return fmt.Errorf("invalid type for key %s: expected %T, got %T",
			c.key, c.val, value)
	}
	c.val = val

	// 发布配置改变事件
	_ = eventCenter.Event.Publish(topic.ConfigChange, &topic.ConfigChangeData{Key: c.key, Value: c.val})
	return nil
}
