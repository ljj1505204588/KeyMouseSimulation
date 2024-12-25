package uiConfig

import (
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
		Label{AssignTo: &c.ifMouseTrackLabel, Text: component.Center.Get(component.MouseTrackStr), ColumnSpan: 4},
		CheckBox{AssignTo: &c.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged: c.setIfTrackMouseMoveClick},
	}
}
func (c *recordConfig) disPlay() []Widget {
	return c.widget
}
func (c *recordConfig) register() {
	component.Center.RegisterChange(c.languageHandler)
	component.RecordConfig.SetMouseTrackChange(false, c.ifTrackMouseMoveRegister)
}

// 语言变动回调
func (c *recordConfig) languageHandler() {
	for c.ifMouseTrackLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	tryPublishErr(c.ifMouseTrackLabel.SetText(component.Center.Get(component.MouseTrackStr)))
}

// 设置是否追踪鼠标移动路径
func (c *recordConfig) setIfTrackMouseMoveClick() {
	component.RecordConfig.SetMouseTrack(c.ifMouseTrackCheck.Checked())
}

// 鼠标移动路径变动回调
func (c *recordConfig) ifTrackMouseMoveRegister(track bool) {
	c.ifMouseTrackCheck.SetChecked(track)
}
