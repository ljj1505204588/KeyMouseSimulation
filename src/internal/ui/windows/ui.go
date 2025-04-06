package uiWindows

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	component_base_button "KeyMouseSimulation/internal/ui/components/baseButton"
	component_config "KeyMouseSimulation/internal/ui/components/config"
	component_menu_items "KeyMouseSimulation/internal/ui/components/menuItems"
	component_system "KeyMouseSimulation/internal/ui/components/system"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	_ "embed"
	"fmt"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
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
		return c.mw.SetTitle(language.MainWindowTitleStr.ToString())
	})

	eventCenter.Event.Register(topic.LanguageChange, func(data interface{}) (err error) {
		t.mw.SetVisible(false)
		t.mw.SetVisible(true)
		return nil
	})

	go eventCenter.Event.Publish(topic.LanguageChange, &topic.LanguageChangeData{
		Typ: enum.LanguageType(conf.LanguageConf.GetValue()),
	})
}

// ----------------------- 主窗口 -----------------------

// //go:embed icon.jpg
// var icon []byte

func MainWindows() {
	// 设置图标
	icon, err := walk.NewIconFromFile(".\\images\\icon_960x960.ico")
	if err != nil {
		fmt.Printf("加载图标失败: %v\n", err)
	}

	var widget []declarative.Widget
	for _, component := range c.widgets {
		widget = append(widget, component.DisPlay(&c.mw)...)
	}
	var menuItems []declarative.MenuItem
	for _, item := range c.menuItems {
		menuItems = append(menuItems, item.MenuItems(&c.mw)...)
	}

	c.Init()

	_, err = declarative.MainWindow{
		AssignTo: c.MWPoint(),
		Title:    language.MainWindowTitleStr.ToString(),
		Size:     declarative.Size{Width: 320, Height: 400},
		Layout:   declarative.Grid{Columns: 8, Alignment: declarative.AlignHNearVCenter},
		Children: widget,
		// 工具栏
		MenuItems: menuItems,
		Icon:      icon,
	}.Run()

	if err != nil {
		fmt.Println(err.Error())
	}
}
