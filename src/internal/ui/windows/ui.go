package uiWindows

import (
	"KeyMouseSimulation/internal/ui"
	uiComponent "KeyMouseSimulation/internal/ui/components"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/event_topic"
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
	eventCenter.Event.Register(event_topic.LanguageChange, func(data interface{}) (err error) {
		t.mw.SetVisible(false)
		t.mw.SetVisible(true)
		return nil
	})
	go eventCenter.Event.Publish(event_topic.LanguageChange, event_topic.LanguageChangeData{})
}

// ----------------------- 主窗口 -----------------------

func MainWindows() (err error) {
	var app = c.MWPoint()
	// todo 设置图标
	var widget []Widget
	for _, component := range c.widgets {
		widget = append(widget, component.DisPlay(&c.mw)...)
	}
	var menuItems []MenuItem
	for _, item := range c.menuItems {
		menuItems = append(menuItems, item.MenuItems(&c.mw)...)
	}

	var mwT = &MainWindow{
		AssignTo: app,
		Title:    language.MainWindowTitleStr.ToString(),
		Size:     Size{Width: 320, Height: 400},
		Layout:   Grid{Columns: 8, Alignment: AlignHNearVCenter},
		Children: widget,
		// 工具栏
		MenuItems: menuItems,
	}
	//c.Init()
	if err = mwT.Create(); err != nil {
		return
	}
	// 手动绑定 Loaded 事件
	mainWindow.Loaded().Attach(func() {
		walk.MsgBox(mainWindow, "提示", "窗口加载完成！", walk.MsgBoxIconInformation)
	})

	// 显示窗口并运行消息循环
	mainWindow.Show()
	walk.RunMainWindow(mainWindow)
	_, err = app.Run()
	return
}
