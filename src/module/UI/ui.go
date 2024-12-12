package ui

import (
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/module/UI/uiComponent"
	"KeyMouseSimulation/module/baseComponent"
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
		&uiComponent.ConfigManageT{},
		&uiComponent.SystemT{},
	},
	menuItems: []KmMenuItem{
		&uiComponent.MenuItemT{},
	},
}

func (t *ControlT) MWPoint() **walk.MainWindow {
	return &t.mw
}

func (t *ControlT) Init() {
	component.Center.RegisterChange(func() {
		t.mw.SetVisible(false)
		t.mw.SetVisible(true)
	})
	go component.Center.SetLanguage(component.Chinese)
}

// ----------------------- 主窗口 -----------------------

func MainWindows() {
	// todo 设置图标
	var widget []Widget
	for _, component := range c.widgets {
		widget = append(widget, component.DisPlay(&c.mw)...)
	}
	var menuItems []MenuItem
	for _, item := range c.menuItems {
		menuItems = append(menuItems, item.MenuItems(&c.mw)...)
	}

	c.Init()
	_, err := MainWindow{
		AssignTo: c.MWPoint(),
		Title:    component.Center.Get(component.MainWindowTitleStr),
		Size:     Size{Width: 320, Height: 400},
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
