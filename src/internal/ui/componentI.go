package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type KmWidget interface {
	DisPlay(mw **walk.MainWindow) []Widget
	LanguageChange(data interface{}) (err error)
}

type KmMenuItem interface {
	MenuItems(mw **walk.MainWindow) []MenuItem
	LanguageChange(data interface{}) (err error)
}
