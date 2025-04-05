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
	"github.com/lxn/walk/declarative"
)

// --------------------------------------- 记录 ----------------------------------------------

type recordConfig struct {

	// 鼠标路径记录
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	recordLenLabel          *walk.Label
	recordLenNumLabel       *walk.NumberLabel
	recordLenLastUpdateTime int64

	widget []declarative.Widget
}

func (c *recordConfig) init() {
	c.widget = []declarative.Widget{
		//鼠标路径
		declarative.Label{AssignTo: &c.ifMouseTrackLabel, Text: language.MouseTrackStr.ToString(), ColumnSpan: 2},
		declarative.CheckBox{AssignTo: &c.ifMouseTrackCheck, ColumnSpan: 2, Checked: true, Alignment: declarative.AlignHCenterVCenter, OnCheckedChanged: func() {
			conf.RecordMouseTrackConf.SetValue(c.ifMouseTrackCheck.Checked())
		}},
		// 记录长度
		declarative.Label{AssignTo: &c.recordLenLabel, Text: language.RecordLenStr.ToString(), ColumnSpan: 2},
		declarative.NumberLabel{AssignTo: &c.recordLenNumLabel, Value: 0, ColumnSpan: 2},
	}

	// 注册回调事件
	eventCenter.Event.Register(topic.ConfigChange, func(data interface{}) (err error) {
		var dataValue = data.(*topic.ConfigChangeData)

		// 监听配置变动
		switch dataValue.Key {
		case enum.RecordMouseTrackConf:
			c.ifMouseTrackCheck.SetChecked(dataValue.Value.(bool))
		case enum.RecordLenConf:
			if nowTime := time.Now().Unix(); nowTime-c.recordLenLastUpdateTime >= 1 {
				_ = c.recordLenNumLabel.SetValue(float64(conf.RecordLen.GetValue()))
				c.recordLenLastUpdateTime = nowTime
			}
		}

		return
	})
}

func (c *recordConfig) disPlay() []declarative.Widget {
	return c.widget
}

// 语言变动回调
func (c *recordConfig) languageChange() {
	for c.ifMouseTrackLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	uiComponent.TryPublishErr(c.ifMouseTrackLabel.SetText(language.MouseTrackStr.ToString()))
}
