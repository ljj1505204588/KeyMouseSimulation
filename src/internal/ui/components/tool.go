package uiComponent

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/topic"
)

func TryPublishErr(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(topic.ServerError, topic.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}
