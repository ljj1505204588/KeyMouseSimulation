package eventCenter

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/share/event_topic"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type factory struct {
	sync.RWMutex
	eventMap map[Topic][]*item
}

var Event EventI = &factory{
	eventMap: make(map[Topic][]*item),
}

type item struct {
	handler Handler // 执行函数
	conf    *config // 配置
}

// Register 注册
func (e *factory) Register(topic Topic, handler Handler, opts ...Options) {
	defer common.RLockSelf(&e.RWMutex)()

	// 默认配置
	var conf = getDefConfig()
	for _, opt := range opts {
		opt(conf)
	}

	// 事件列表
	var list = e.eventMap[topic]
	list = append(list, &item{
		handler: handler,
		conf:    conf,
	})

	sort.Slice(list, func(i, j int) bool {
		return list[i].conf.order > list[j].conf.order
	})
	e.eventMap[topic] = list
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
func (e *factory) getHandler(topic Topic) (list []Handler, ok bool) {
	defer common.RRLockSelf(&e.RWMutex)()

	var items []*item
	if items, ok = e.eventMap[topic]; ok {
		for _, perItem := range items {
			list = append(list, perItem.handler)
		}
	}

	return
}
