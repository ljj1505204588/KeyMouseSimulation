package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type KmWidget interface {
	DisPlay(mw **walk.MainWindow) []Widget
}

type KmMenuItem interface {
	MenuItems(mw **walk.MainWindow) []MenuItem
}
