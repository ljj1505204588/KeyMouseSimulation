package uiComponent

import (
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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
		conf.register()
	}
}

func (t *ConfigManageT) DisPlay(mw **walk.MainWindow) (res []Widget) {
	t.mw = mw
	t.Once.Do(t.Init)

	for _, conf := range t.configs {
		res = append(res, conf.disPlay()...)
	}
	return
}
