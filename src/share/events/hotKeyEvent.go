package events

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
)

// ServerHotKeyDown 热键按下
const ServerHotKeyDown eventCenter.Topic = "server_hot_key_down"

type ServerHotKeyDownData struct {
	HotKey enum.HotKey
	Key    *string
}
