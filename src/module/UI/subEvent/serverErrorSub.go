package subEvent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	ui "KeyMouseSimulation/module/UI"
	"KeyMouseSimulation/share/events"
)

type ServerErrorSub struct {
	c ui.ControlT
}

func NewServerErrorSub(c ui.ControlT) {
	topic := events.ServerError
	handler := ServerErrorSub{c: c}
	eventCenter.Event.Register(topic, handler)
}

func (sub ServerErrorSub) Handler(data interface{}) (err error) {
	subData := data.(events.ServerErrorData)

	return sub.c.ShowError(subData.ErrInfo)
}
