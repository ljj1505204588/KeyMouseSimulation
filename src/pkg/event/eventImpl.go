package eventCenter

import "KeyMouseSimulation/share/topic"

type EventI interface {
	Register(topic topic.Topic, handler Handler, opts ...Options) // 注册
	Publish(topic topic.Topic, dataAny interface{}) (err error)   // 同步发布
	ASyncPublish(topic topic.Topic, dataAny interface{})          // 异步发布
}

type Handler func(data interface{}) (err error)
