package uiComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/events"
)

type BaseT struct {
}

func tryPublishErr(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
	}
}
