package BaseComponent

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

	//文件选择
	basePath  string
	fileBox   *walk.ComboBox
	fileNames []string

	//系统状态
	Status int
}

func (b *BaseT) Init() {
	b.hKList = [4]string{"F7", "F8", "F9", "F10"}
}

// --------------------------------------- 基础功能 ----------------------------------------------

func (b *BaseT) initCheck() bool {
	return b.baseInitCheck()
}

// 初始化校验
func (b *BaseT) baseInitCheck() bool {

	return b.mw != nil && b.fileBox != nil
}

func (b *BaseT) lockSelf() func() {
	b.l.Lock()
	return b.l.Unlock
}
