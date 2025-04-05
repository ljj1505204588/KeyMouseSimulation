package component_config

import (
	"sync"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type ConfigManageT struct {
	mw **walk.MainWindow
	sync.Once

	configs []configI // 配置
}

func (t *ConfigManageT) Init() {
	t.configs = []configI{&fileConfig{}, &recordConfig{}, &playbackConfig{}}
	for _, conf := range t.configs {
		conf.init()
	}
}

func (t *ConfigManageT) DisPlay(mw **walk.MainWindow) (res []declarative.Widget) {
	t.mw = mw
	t.Once.Do(t.Init)

	for _, conf := range t.configs {
		res = append(res, conf.disPlay()...)
	}
	return
}

// LanguageChange 设置语言
func (t *ConfigManageT) LanguageChange(data interface{}) (err error) {
	for _, conf := range t.configs {
		conf.languageChange()
	}

	return
}
