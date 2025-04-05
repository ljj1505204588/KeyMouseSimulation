package component_config

import (
	"github.com/lxn/walk/declarative"
)

type configI interface {
	init()                         // 初始化
	disPlay() []declarative.Widget // 展示
	languageChange()               // 语言设置
}
