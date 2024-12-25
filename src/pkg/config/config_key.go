// config_key.go
package conf

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/event_topic"
	"fmt"
)

type ConfigKeyI interface {
	GetKey() string
	GetValue() any
	SetValue(value any) error
}

// configKeyT 泛型配置键
type configKeyT[T any] struct {
	key string
	val T
}

func (c *configKeyT[T]) GetKey() string {
	return c.key
}

func (c *configKeyT[T]) GetValue() any {
	return c.val
}
func (c *configKeyT[T]) SetValue(value any) error {
	val, ok := value.(T)
	if !ok {
		return fmt.Errorf("invalid type for key %s: expected %T, got %T",
			c.key, c.val, value)
	}
	c.val = val

	// 发布配置改变事件
	eventCenter.Event.Publish(event_topic.ConfigChange, event_topic.ConfigChangeData{Key: c.key})
	return nil
}

// 预定义配置键
var (
	KeyRecordHotKey     = &configKeyT[string]{key: "record.hotkey", val: "F9"}
	KeyPlaybackHotKey   = &configKeyT[string]{key: "playback.hotkey", val: "F10"}
	KeyPauseHotKey      = &configKeyT[string]{key: "pause.hotkey", val: "F11"}
	KeyStopHotKey       = &configKeyT[string]{key: "stop.hotkey", val: "F12"}
	KeyRecordMouseTrack = &configKeyT[bool]{key: "record.mousetrack", val: true}
	KeyPlaybackSpeed    = &configKeyT[float64]{key: "playback.speed", val: 1.0}
	KeyPlaybackTimes    = &configKeyT[int64]{key: "playback.times", val: int64(1)}
	KeyLanguage         = &configKeyT[string]{key: "system.language", val: "zh_CN"}
)

// 获取所有配置键
func allKeys() []ConfigKeyI {
	return []ConfigKeyI{
		KeyRecordHotKey,
		KeyPlaybackHotKey,
		KeyPauseHotKey,
		KeyStopHotKey,
		KeyRecordMouseTrack,
		KeyPlaybackSpeed,
		KeyPlaybackTimes,
		KeyLanguage,
	}
}
