package uiConfig

import (
	"KeyMouseSimulation/internal/core/file"
	"KeyMouseSimulation/pkg/common/gene"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/event_topic"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// --------------------------------------- 文件 ----------------------------------------------

type configI interface {
	init()
	disPlay() []Widget
	register()
}
type fileConfig struct {
	fileNames     []string
	fileComponent file.FileControlI

	fileLabel *walk.Label
	fileBox   *walk.ComboBox

	widget []Widget
}

func (c *fileConfig) init() {
	c.widget = []Widget{
		Label{AssignTo: &c.fileLabel, ColumnSpan: 2},
		ComboBox{AssignTo: &c.fileBox, ColumnSpan: 6, Model: c.fileNames, Editable: true, OnCurrentIndexChanged: c.chooseFile},
	}

	c.fileComponent = file.FileControl
}
func (c *fileConfig) disPlay() []Widget {

	return c.widget
}

func (c *fileConfig) register() {
	eventCenter.Event.Register(event_topic.LanguageChange, c.languageHandler)
	c.fileComponent.FileChange(c.fileChangeHandler)
}

// 语言变动回调
func (c *fileConfig) languageHandler(data interface{}) (err error) {
	for c.fileLabel == nil {
		time.Sleep(20 * time.Millisecond)
	}

	// 设置文件标签
	if err = c.fileLabel.SetText(language.FileLabelStr.ToString()); err != nil {
		eventCenter.Event.Publish(event_topic.ServerError, event_topic.ServerErrorData{ErrInfo: err.Error()})
	}
	return
}

// chooseFile 选择文件
func (c *fileConfig) chooseFile() {
	tryPublishErr(c.fileComponent.Choose(c.fileBox.Text()))
}

// 文件变动
func (c *fileConfig) fileChangeHandler(names []string, newFile []string) {
	if c.fileLabel == nil || c.fileBox == nil {
		return
	}

	// 模式设置
	tryPublishErr(c.fileBox.SetModel(names))

	// 新文件设置
	if len(newFile) != 0 {
		tryPublishErr(c.fileComponent.Choose(newFile[0]))
		tryPublishErr(c.fileBox.SetText(newFile[0]))
	}

	// 文件删除重设
	if !gene.Contain(names, c.fileBox.Text()) {
		var current = ""
		if len(names) > 0 {
			current = names[0]
		}
		tryPublishErr(c.fileComponent.Choose(current))
		tryPublishErr(c.fileBox.SetText(current))
	}

	return
}

func tryPublishErr(err error) {
	if err != nil {
		eventCenter.Event.Publish(event_topic.ServerError, event_topic.ServerErrorData{ErrInfo: err.Error()})
	}
}
