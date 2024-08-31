package eventCenter

import "errors"

type factory struct {
	eventMap map[Topic][]Handler
}

var Event = factory{
	make(map[Topic][]Handler),
}

// Register 注册
func (e *factory) Register(topic Topic, handler Handler) {
	e.eventMap[topic] = append(e.eventMap[topic], handler)
}

// Publish 同步
func (e *factory) Publish(topic Topic, data interface{}) (err error) {
	handlers, ok := e.eventMap[topic]
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
func (e *factory) ASyncPublish(topic Topic, data interface{}) (err error) {
	handlers, ok := e.eventMap[topic]
	if !ok {
		return errors.New("Topic Unregistered. ")

	}

	for _, h := range handlers {
		go h(data)
	}

	return
}
