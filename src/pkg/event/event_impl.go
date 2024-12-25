package eventCenter

type EventI interface {
	Register(topic Topic, handler Handler)             // 注册
	Publish(topic Topic, data interface{}) (err error) // 同步发布
	ASyncPublish(topic Topic, data interface{})        // 异步发布
}

type Topic string

type Handler func(data interface{}) (err error)
