package component_config

import (
	uiComponent "KeyMouseSimulation/internal/ui/components"
	eventCenter "KeyMouseSimulation/pkg/event"
	rp_file "KeyMouseSimulation/pkg/file"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/topic"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

// --------------------------------------- 文件 ----------------------------------------------

type fileConfig struct {
	fileLabel *walk.Label
	fileBox   *walk.ComboBox

	widget []declarative.Widget
}

func (c *fileConfig) init() {
	var name, fileList = rp_file.FileControl.Current()

	c.widget = []declarative.Widget{
		declarative.Label{AssignTo: &c.fileLabel, ColumnSpan: 2},
		declarative.ComboBox{AssignTo: &c.fileBox, ColumnSpan: 6, Value: name, Model: fileList, Editable: true, OnCurrentIndexChanged: c.chooseFile},
	}

	c.register()
}
func (c *fileConfig) disPlay() []declarative.Widget {
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

	for c.fileLabel == nil || c.fileBox == nil {
		time.Sleep(50 * time.Millisecond)
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
