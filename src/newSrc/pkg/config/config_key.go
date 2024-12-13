// config_key.go
package conf

import (
	"fmt"
)

type ConfigKeyI interface {
	getKey() string
	getDefault() (any, bool)
	validate(value any) error
}

// configKeyT 泛型配置键
type configKeyT[T any] struct {
	key string
	val T
}

func (c *configKeyT[T]) getKey() string {
	return c.key
}

func (c *configKeyT[T]) getDefault() (any, bool) {
	return c.val, true
}

func (c *configKeyT[T]) validate(value any) error {
	_, ok := value.(T)
	if !ok {
		return fmt.Errorf("invalid type for key %s: expected %T, got %T",
			c.key, c.val, value)
	}
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
