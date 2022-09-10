package subEvent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	ui "KeyMouseSimulation/module/UI"
	"KeyMouseSimulation/share/events"
)

type ServerStatusChangeSub struct {
	c ui.ControlT
}

func NewServerStatusChangeSub(c ui.ControlT) {
	topic := events.ServerStatusChange
	handler := ServerStatusChangeSub{c: c}
	eventCenter.Event.Register(topic, handler)
}

func (sub ServerStatusChangeSub) Handler(data interface{}) (err error) {
	subData := data.(events.ServerStatusChangeData)

	return sub.c.StageChange(subData.Status)
}
