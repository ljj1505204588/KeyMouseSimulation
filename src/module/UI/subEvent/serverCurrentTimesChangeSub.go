package subEvent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	ui "KeyMouseSimulation/module/UI"
	"KeyMouseSimulation/share/events"
)

type ServerCurrentTimesChangeSub struct {
	c ui.ControlT
}

func NewServerCurrentTimesChangeSub(c ui.ControlT) {
	topic := events.ServerCurrentTimesChange
	handler := ServerCurrentTimesChangeSub{c: c}
	eventCenter.Event.Register(topic, handler)
}

func (sub ServerCurrentTimesChangeSub) Handler(data interface{}) (err error) {
	subData := data.(events.ServerCurrentTimesChangeData)

	return sub.c.SetCurrentTimes(subData.CurrentTimes)
}
