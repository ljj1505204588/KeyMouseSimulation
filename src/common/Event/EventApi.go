package eventCenter

type handler interface {
	Handler(data interface{}) (err error)
}
