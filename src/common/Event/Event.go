package eventCenter

type Handler func(data interface{}) (err error)
