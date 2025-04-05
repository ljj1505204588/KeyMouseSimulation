package eventCenter

type Options func(opt *config)

type config struct {
	order int
}

func getDefConfig() *config {
	return &config{order: -1}
}

func SetOrderLv(lv int) Options {
	return func(opt *config) {
		opt.order = lv
	}
}
