package uiComponent

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type KmWidget interface {
	DisPlay(mw **walk.MainWindow) []declarative.Widget
	LanguageChange(data interface{}) (err error)
}

type KmMenuItem interface {
	MenuItems(mw **walk.MainWindow) []declarative.MenuItem
	LanguageChange(data interface{}) (err error)
}
