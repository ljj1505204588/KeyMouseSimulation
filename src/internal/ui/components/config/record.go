package component_config

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// --------------------------------------- 记录 ----------------------------------------------

type recordConfig struct {

	// 鼠标路径记录
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	widget []Widget
}

func (c *recordConfig) init() {
	c.widget = []Widget{
		//鼠标路径
		Label{AssignTo: &c.ifMouseTrackLabel, Text: language.MouseTrackStr.ToString(), ColumnSpan: 4},
		CheckBox{AssignTo: &c.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged: func() {
			uiComponent.TryPublishErr(conf.RecordMouseTrackConf.SetValue(c.ifMouseTrackCheck.Checked()))
		}},
	}

	// 注册回调事件
	eventCenter.Event.Register(topic.ConfigChange, func(data interface{}) (err error) {
		var dataValue = data.(*topic.ConfigChangeData)

		// 监听配置变动
		if dataValue.Key == enum.RecordMouseTrackConf {
			c.ifMouseTrackCheck.SetChecked(dataValue.Value.(bool))
		}

		return
	})
}

func (c *recordConfig) disPlay() []Widget {
	return c.widget
}

// 语言变动回调
func (c *recordConfig) languageChange() {
	for c.ifMouseTrackLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	uiComponent.TryPublishErr(c.ifMouseTrackLabel.SetText(language.MouseTrackStr.ToString()))
}
