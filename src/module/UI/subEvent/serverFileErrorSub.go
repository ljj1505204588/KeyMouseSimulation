package subEvent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	ui "KeyMouseSimulation/module/UI"
	"KeyMouseSimulation/share/events"
)

type ServerFileErrorSub struct {
	c ui.ControlT
}

func NewServerFileErrorSub(c ui.ControlT) {
	topic := events.ServerFileError
	handler := ServerFileErrorSub{c: c}
	eventCenter.Event.Register(topic, handler)
}

func (sub ServerFileErrorSub) Handler(data interface{}) (err error) {
	subData := data.(events.ServerFileErrorData)

	return sub.c.ShowFileError(subData.ErrInfo)
}
