package paramoption

type Option func(o *OptionT)

type OptionT struct {
	fromCent int
}

func SetCent() func(o *OptionT) {
	return func(o *OptionT) {
		o.fromCent = 1
	}
}
