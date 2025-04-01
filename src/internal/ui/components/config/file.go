package component_config

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/file"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/topic"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// --------------------------------------- 文件 ----------------------------------------------

type fileConfig struct {
	fileNames []string

	fileLabel *walk.Label
	fileBox   *walk.ComboBox

	widget []Widget
}

func (c *fileConfig) init() {
	c.widget = []Widget{
		Label{AssignTo: &c.fileLabel, ColumnSpan: 2},
		ComboBox{AssignTo: &c.fileBox, ColumnSpan: 6, Model: c.fileNames, Editable: true, OnCurrentIndexChanged: c.chooseFile},
	}

}
func (c *fileConfig) disPlay() []Widget {

	return c.widget
}

func (c *fileConfig) register() {
	eventCenter.Event.Register(topic.FileListChange, c.fileChangeHandler)
}

// chooseFile 选择文件
func (c *fileConfig) chooseFile() {
	uiComponent.TryPublishErr(rp_file.FileControl.Choose(c.fileBox.Text()))
}

// 文件变动
func (c *fileConfig) fileChangeHandler(data interface{}) (err error) {
	var dataValue = data.(*topic.FileListChangeData)

	if c.fileLabel == nil || c.fileBox == nil {
		return
	}

	// 模式设置
	uiComponent.TryPublishErr(c.fileBox.SetModel(dataValue.Files))
	uiComponent.TryPublishErr(c.fileBox.SetText(dataValue.ChooseFile))

	return
}

// 语言变动回调
func (c *fileConfig) languageChange() {
	for c.fileLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	// 设置文件标签
	uiComponent.TryPublishErr(c.fileLabel.SetText(language.FileLabelStr.ToString()))
}
