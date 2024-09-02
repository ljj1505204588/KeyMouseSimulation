package svcComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/events"
)

// 发布服务错误事件
func tryPublishServerError(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
	}
}

// 发布服务错误事件
func publishServerError(err error) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
}
