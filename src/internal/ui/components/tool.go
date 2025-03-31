package uiComponent

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/event_topic"
)

func tryPublishErr(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(event_topic.ServerError, event_topic.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}
