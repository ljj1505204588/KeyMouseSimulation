package uiWindows

import (
	"KeyMouseSimulation/internal/ui"
	uiComponent "KeyMouseSimulation/internal/ui/components"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type ControlT struct {
	mw *walk.MainWindow

	widgets   []ui.KmWidget
	menuItems []ui.KmMenuItem
}

var c = &ControlT{
	widgets: []ui.KmWidget{
		&uiComponent.FunctionT{},     // 功能
		&uiComponent.ConfigManageT{}, // 配置
		&uiComponent.SystemT{},       // 系统
	},
	menuItems: []ui.KmMenuItem{
		&uiComponent.MenuItemT{},
	},
}

func (t *ControlT) MWPoint() **walk.MainWindow {
	return &t.mw
}

func (t *ControlT) Init() {
	eventCenter.Event.Register(eventCenter.LanguageChange, func() {
		t.mw.SetVisible(false)
		t.mw.SetVisible(true)
	})
	go eventCenter.Event.Publish(eventCenter.LanguageChange, eventCenter.LanguageChangeData{})
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
