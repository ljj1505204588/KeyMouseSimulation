package BaseComponent

import (
	"KeyMouseSimulation/module/server"
	"KeyMouseSimulation/module2/language"
	"github.com/lxn/walk"
	"sync"
	"time"
)

type BaseT struct {
	sc server.ControlI
	mw *walk.MainWindow
	l  sync.Mutex

	languageMap     map[language.DisPlay]string
	languageTyp     language.LanguageTyp
	languageHandler []func(typ language.LanguageTyp)

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

func (b *BaseT) Init(sc server.ControlI) {
	b.sc = sc

	b.hKList, b.keyList = b.sc.GetKeyList()
}

func (b *BaseT) ChangeLanguage(typ language.LanguageTyp, sync bool) {
	var f = func() {
		var m = language.LanguageMap[typ]
		b.languageMap = language.LanguageMap[typ]
		b.languageTyp = typ

		for _, h := range b.languageHandler {
			h(typ)
		}

		b.waitInitCheck()

		_ = b.mw.SetTitle(m[language.MainWindowTitleStr])

		b.mw.SetVisible(false)
		b.mw.SetVisible(true)
	}

	if sync {
		go f()
	} else {
		f()
	}
}

// --------------------------------------- 基础功能 ----------------------------------------------

func (b *BaseT) waitInitCheck() {
	for !b.baseInitCheck() {
		time.Sleep(10 * time.Millisecond)
	}
}

// 初始化校验
func (b *BaseT) baseInitCheck() bool {
	//for _, per := range b.hKBox {
	//	if per == nil {
	//		return false
	//	}
	//}

	return b.mw != nil && b.fileBox != nil
}

// 注册修改语言函数
func (b *BaseT) registerChangeLanguage(h ...func(typ language.LanguageTyp)) {
	b.languageHandler = append(b.languageHandler, h...)
}

func (b *BaseT) lockSelf() func() {
	b.l.Lock()
	return b.l.Unlock
}
