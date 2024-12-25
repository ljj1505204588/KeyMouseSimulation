package eventCenter

import (
	"KeyMouseSimulation/common/commonTool"
	"KeyMouseSimulation/share/event_topic"
	"errors"
	"fmt"
	"sync"
)

type factory struct {
	sync.RWMutex
	eventMap map[Topic][]Handler
}

var Event EventI = &factory{
	eventMap: make(map[Topic][]Handler),
}

// Register 注册
func (e *factory) Register(topic Topic, handler Handler) {
	defer commonTool.RLockSelf(&e.RWMutex)()
	e.eventMap[topic] = append(e.eventMap[topic], handler)
}

// Publish 同步
func (e *factory) Publish(topic Topic, data interface{}) (err error) {
	handlers, ok := e.getHandler(topic)
	if !ok {
		return errors.New("Topic Unregistered. ")
	}

	for _, h := range handlers {
		if err = h(data); err != nil {
			return
		}
	}

	return
}

// ASyncPublish 异步
func (e *factory) ASyncPublish(topic Topic, data interface{}) {
	handlers, ok := e.getHandler(topic)
	if !ok {
		return
	}

	for _, h := range handlers {
		go func(h Handler) {
			if err := h(data); err != nil {
				var errInfo = fmt.Sprintf("异步执行事件[%s]失败, 错误信息: %s", topic, err.Error())
				_ = e.Publish(event_topic.ServerError, event_topic.ServerErrorData{ErrInfo: errInfo})
			}
		}(h)
	}

	return
}
func (e *factory) getHandler(topic Topic) ([]Handler, bool) {
	defer commonTool.RRLockSelf(&e.RWMutex)()
	handlers, ok := e.eventMap[topic]
	return handlers, ok
}
