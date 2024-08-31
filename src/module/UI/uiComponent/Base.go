package uiComponent

import (
	"github.com/lxn/walk"
	"sync"
)

type BaseT struct {
	l sync.Mutex

	//热键
	keyList []string
	hKList  [4]string
	hKBox   [4]*walk.ComboBox
}

func (b *BaseT) Init() {
	b.hKList = [4]string{"F7", "F8", "F9", "F10"}
}
