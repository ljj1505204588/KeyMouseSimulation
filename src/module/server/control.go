package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"sync"
)

var control controlT

func Init() {

	control = controlT{
		kmStatusI: chooseBox[enum.Free],
	}
	eventCenter.Event.Register(events.ButtonClick, (&control).buttonClickHandler)
}

type controlT struct {
	l sync.Mutex
	kmStatusI
}

// Record 记录
func (c *controlT) Record() {
	defer c.lockSelf()()

	c.kmStatusI.record()
	c.kmStatusI = c.kmStatusI.current()
}

// Playback 回放
func (c *controlT) Playback(name string) {
	defer c.lockSelf()()

	c.kmStatusI.playback(name)
	c.kmStatusI = c.kmStatusI.current()
}

// Pause 暂停
func (c *controlT) Pause() {
	defer c.lockSelf()()

	c.kmStatusI.pause()
	c.kmStatusI = c.kmStatusI.current()
}

// Stop 停止
func (c *controlT) Stop() {
	defer c.lockSelf()()

	c.kmStatusI.stop()
	c.kmStatusI = c.kmStatusI.current()
}

// Save 存储
func (c *controlT) Save(name string) {
	//defer c.lockSelf()()

	c.kmStatusI = c.kmStatusI.current()
	c.kmStatusI.save(name)
	c.kmStatusI = c.kmStatusI.current()
}

func (c *controlT) lockSelf() func() {
	c.l.Lock()
	return func() {
		c.l.Unlock()
		c.publishStatusChange()
	}
}
func (c *controlT) publishStatusChange() {
	_ = eventCenter.Event.Publish(events.ServerStatusChange, events.ServerStatusChangeData{
		Status: c.status(),
	})
}

func (c *controlT) buttonClickHandler(data interface{}) (err error) {
	var d = data.(events.ButtonClickData)
	switch d.Button {
	case enum.RecordButton:
		c.Record()
	case enum.PlaybackButton:
		c.Playback(d.Name)
	case enum.PauseButton:
		c.Pause()
	case enum.StopButton:
		c.Stop()
	case enum.SaveFileButton:
		c.Save(d.Name)
	}
	return
}
