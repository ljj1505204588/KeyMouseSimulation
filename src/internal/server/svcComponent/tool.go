package svcComponent

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/topic"
)

// 发布服务错误事件
func tryPublishServerError(err error) {
	if err != nil {
		eventCenter.Event.ASyncPublish(topic.ServerError, topic.ServerErrorData{ErrInfo: err.Error()})
	}
}

// 发布服务错误事件
func publishServerError(err error) {
	eventCenter.Event.ASyncPublish(topic.ServerError, topic.ServerErrorData{ErrInfo: err.Error()})
}
