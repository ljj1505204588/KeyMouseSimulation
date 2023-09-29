package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/module/server/BaseComponent"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"sync"
)

var control controlT

func Init() {
	initKmStatus()

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
	defer c.lockSelf()()

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

// ----------------------------------- 状态机 -----------------------------------

type kmStatusI interface {
	record()              // 记录
	playback(name string) // 回放
	pause()               // 暂停
	stop()                // 停止
	save(name string)     // 存储

	current() kmStatusI
	status() enum.Status
}

var chooseBox map[enum.Status]kmStatusI

func initKmStatus() {
	var base = &BaseT{
		Status:   enum.Free,
		PlayBack: recordAndPlayBack.GetPlaybackServer(),
		Record:   recordAndPlayBack.GetRecordServer(),
	}
	chooseBox = map[enum.Status]kmStatusI{
		enum.Free:          &freeStatusT{&baseStatusT{base: base}},
		enum.Recording:     &recordingStatusT{&baseStatusT{base: base}},
		enum.RecordPause:   &recordPauseStatusT{&baseStatusT{base: base}},
		enum.Playback:      &playbackStatusT{&baseStatusT{base: base}},
		enum.PlaybackPause: &playbackPauseStatusT{&baseStatusT{base: base}},
	}
}

type BaseT struct {
	Status   enum.Status
	PlayBack recordAndPlayBack.PlayBackServerI
	Record   recordAndPlayBack.RecordServerI
}

// --------------------------------- 状态机 ---------------------------------

// "Free" 状态机
type freeStatusT struct {
	*baseStatusT
}

func (s *freeStatusT) record() {
	s.base.Record.Start()
	s.base.Status = enum.Recording
}
func (s *freeStatusT) playback(name string) {
	s.base.PlayBack.Start(name)
	s.base.Status = enum.Playback
}
func (s *freeStatusT) save(name string) {
	s.base.Record.Save(name)
}
func (s *freeStatusT) status() enum.Status {
	return enum.Free
}

// "Recording" 状态机
type recordingStatusT struct {
	*baseStatusT
}

func (s *recordingStatusT) pause() {
	s.base.Record.Pause()
	s.base.Status = enum.RecordPause
}
func (s *recordingStatusT) stop() {
	s.base.Record.Stop()
	s.base.Status = enum.Free
	_ = eventCenter.Event.Publish(events.RecordFinish, events.RecordFinishData{})
}
func (s *recordingStatusT) status() enum.Status {
	return enum.Recording
}

// "RecordPause" 状态机
type recordPauseStatusT struct {
	*baseStatusT
}

func (s *recordPauseStatusT) record() {
	s.base.Record.Start()
	s.base.Status = enum.Recording
}
func (s *recordPauseStatusT) stop() {
	s.base.Record.Stop()
	s.base.Status = enum.Free
}
func (s *recordPauseStatusT) status() enum.Status {
	return enum.RecordPause
}

// "Playback" 状态机
type playbackStatusT struct {
	*baseStatusT
}

func (s *playbackStatusT) pause() {
	s.base.PlayBack.Pause()
	s.base.Status = enum.PlaybackPause
}
func (s *playbackStatusT) stop() {
	s.base.PlayBack.Stop()
	s.base.Status = enum.Free
}
func (s *playbackStatusT) status() enum.Status {
	return enum.Playback
}

// "PlaybackPause" 状态机
type playbackPauseStatusT struct {
	*baseStatusT
}

func (s *playbackPauseStatusT) playback(name string) {
	s.base.PlayBack.Start(name)
	s.base.Status = enum.Playback
}
func (s *playbackPauseStatusT) stop() {
	s.base.PlayBack.Stop()
	s.base.Status = enum.Free
}
func (s *playbackPauseStatusT) status() enum.Status {
	return enum.PlaybackPause
}

// --------------------------------- 基础状态 ---------------------------------

type baseStatusT struct {
	base *BaseT
}

func (s *baseStatusT) record() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) playback(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) pause() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) stop() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) save(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}

func (s *baseStatusT) current() kmStatusI {
	return chooseBox[s.base.Status]
}

func (s *baseStatusT) status() enum.Status {
	return s.base.Status
}
