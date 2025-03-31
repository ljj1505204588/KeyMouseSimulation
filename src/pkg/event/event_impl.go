package eventCenter

type EventI interface {
	Register(topic Topic, handler Handler, opts ...Options) // 注册
	Publish(topic Topic, dataAny interface{}) (err error)   // 同步发布
	ASyncPublish(topic Topic, dataAny interface{})          // 异步发布
}

type Topic string

type Handler func(data interface{}) (err error)
