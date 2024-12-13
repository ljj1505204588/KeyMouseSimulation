package eventCenter

type Topic string

type Handler func(data interface{}) (err error)
