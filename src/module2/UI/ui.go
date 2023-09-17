package ui

import (
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/module/server"
	"KeyMouseSimulation/module2/UI/BaseComponent"
	language2 "KeyMouseSimulation/module2/language"
	"KeyMouseSimulation/share/language"
	. "github.com/lxn/walk/declarative"
	"time"
)

type ControlT struct {
	wc server.ControlI

	function BaseComponent.FunctionT
	playBack BaseComponent.PlaybackT
	system   BaseComponent.SystemT
}

var c *ControlT

func createControl() *ControlT {
	c = &ControlT{}
	c.wc = server.NewWinControl()

	var base = BaseComponent.BaseT{}
	base.Init(c.wc)

	c.function.Init(&base)
	c.system.Init(&base)
	c.playBack.Init(&base)

	base.ChangeLanguage(language2.Chinese, true)
	return c
}

// ----------------------- 主窗口 -----------------------

func MainWindows() {
	c = createControl()
	var widget []Widget
	widget = append(widget, c.function.DisPlay()...)
	widget = append(widget, c.playBack.DisPlay()...)
	widget = append(widget, c.system.DisPlay()...)

	_, err := MainWindow{
		AssignTo: c.system.MWPoint(),
		Title:    language.MainWindowTitleStr,
		Size:     Size{Width: 320, Height: 240},
		Layout:   Grid{Columns: 8, Alignment: AlignHNearVCenter},
		Children: widget,
		//工具栏
		MenuItems: c.system.MenuItems(),
	}.Run()

	if err != nil {
		logTool.ErrorAJ(err)
		time.Sleep(5 * time.Second)
	}
}
