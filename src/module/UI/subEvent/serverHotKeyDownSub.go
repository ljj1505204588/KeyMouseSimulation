package subEvent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	ui "KeyMouseSimulation/module/UI"
	"KeyMouseSimulation/share/events"
)

type ServerHotKeyDownSub struct {
	c ui.ControlT
}

func NewServerHotKeyDownSub(c ui.ControlT) {
	topic := events.ServerHotKeyDown
	handler := ServerHotKeyDownSub{c: c}
	eventCenter.Event.Register(topic, handler)
}

func (sub ServerHotKeyDownSub) Handler(data interface{}) (err error) {
	subData := data.(events.ServerHotKeyDownData)

	return sub.c.HotKeyDown(subData.HotKey)
}
