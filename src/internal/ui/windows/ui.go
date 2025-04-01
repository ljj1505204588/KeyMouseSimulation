//go:build windows
// +build windows

package uiWindows

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	component_base_button "KeyMouseSimulation/internal/ui/components/base_button"
	component_config "KeyMouseSimulation/internal/ui/components/config"
	component_menu_items "KeyMouseSimulation/internal/ui/components/menu_items"
	component_system "KeyMouseSimulation/internal/ui/components/system"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/topic"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type ControlT struct {
	mw *walk.MainWindow

	widgets   []uiComponent.KmWidget
	menuItems []uiComponent.KmMenuItem
}

var c = &ControlT{
	widgets: []uiComponent.KmWidget{
		&component_base_button.FunctionT{}, // 功能
		&component_config.ConfigManageT{},  // 配置
		&component_system.SystemT{},        // 系统
	},
	menuItems: []uiComponent.KmMenuItem{
		&component_menu_items.MenuItemT{},
	},
}

func (t *ControlT) MWPoint() **walk.MainWindow {
	return &t.mw
}

func (t *ControlT) Init() {
	for _, widget := range t.widgets {
		eventCenter.Event.Register(topic.LanguageChange, widget.LanguageChange)
	}
	for _, item := range t.menuItems {
		eventCenter.Event.Register(topic.LanguageChange, item.LanguageChange)
	}

	eventCenter.Event.Register(topic.LanguageChange, func(data interface{}) (err error) {
		t.mw.SetVisible(false)
		t.mw.SetVisible(true)
		return nil
	})

	go eventCenter.Event.Publish(topic.LanguageChange, &topic.LanguageChangeData{})
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
		Title:    language.MainWindowTitleStr.ToString(),
		Size:     Size{Width: 320, Height: 400},
		Layout:   Grid{Columns: 8, Alignment: AlignHNearVCenter},
		Children: widget,
		// 工具栏
		MenuItems: menuItems,
	}.Run()

	if err != nil {
		fmt.Println(err.Error())
	}
}
