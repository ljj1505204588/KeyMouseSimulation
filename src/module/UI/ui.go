package ui

import (
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/module/UI/uiComponent"
	"KeyMouseSimulation/module/language"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

type ControlT struct {
	mw *walk.MainWindow

	widgets   []KmWidget
	menuItems []KmMenuItem
}

var c = &ControlT{
	widgets: []KmWidget{
		&uiComponent.FunctionT{},
		&uiComponent.ConfigT{},
		&uiComponent.SystemT{},
	},
	menuItems: []KmMenuItem{
		&uiComponent.MenuItemT{},
	},
}

func (t *ControlT) MWPoint() **walk.MainWindow {
	return &t.mw
}

// ----------------------- 主窗口 -----------------------

func MainWindows() {
	// todo 设置图标
	var widget []Widget
	for _, component := range c.widgets {
		widget = append(widget, component.DisPlay(c.mw)...)
	}
	var menuItems []MenuItem
	for _, item := range c.menuItems {
		menuItems = append(menuItems, item.MenuItems(c.mw)...)
	}

	_, err := MainWindow{
		AssignTo: c.MWPoint(),
		Title:    language.Center.Get(language.MainWindowTitleStr),
		Size:     Size{Width: 320, Height: 240},
		Layout:   Grid{Columns: 8, Alignment: AlignHNearVCenter},
		Children: widget,
		//工具栏
		MenuItems: menuItems,
	}.Run()

	if err != nil {
		logTool.ErrorAJ(err)
		time.Sleep(5 * time.Second)
	}
}
