package events

import (
	eventCenter "KeyMouseSimulation/common/Event"
)

// ServerHotKeyDown 热键按下
const ServerHotKeyDown eventCenter.Topic = "server_hot_key_down"

type ServerHotKeyDownData struct {
	Key string
}

// SetHotKey 设置热键
const SetHotKey eventCenter.Topic = "set_hot_key"

type SetHotKeyData struct {
	Key string
}
