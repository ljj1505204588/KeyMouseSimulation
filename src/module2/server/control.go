package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module2/language"
	"KeyMouseSimulation/module2/server/BaseComponent"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"sync"
)

func NewControl() ControlI {
	return &ControlT{
		kmStatusI: chooseBox[enum.Free],
	}
}

type ControlT struct {
	l sync.Mutex
	kmStatusI
}

// Record 记录
func (c *ControlT) Record() {
	defer c.lockSelf()()

	c.kmStatusI.record()
	c.kmStatusI = c.kmStatusI.current()
}

// Playback 回放
func (c *ControlT) Playback(name string) {
	defer c.lockSelf()()

	c.kmStatusI.playback(name)
	c.kmStatusI = c.kmStatusI.current()
}

// Pause 暂停
func (c *ControlT) Pause() {
	defer c.lockSelf()()

	c.kmStatusI.pause()
	c.kmStatusI = c.kmStatusI.current()
}

// Stop 停止
func (c *ControlT) Stop() {
	defer c.lockSelf()()

	c.kmStatusI.stop()
	c.kmStatusI = c.kmStatusI.current()
}

// Save 存储
func (c *ControlT) Save(name string) {
	defer c.lockSelf()()

	c.kmStatusI.save(name)
	c.kmStatusI = c.kmStatusI.current()
}

func (c *ControlT) lockSelf() func() {
	c.l.Lock()
	return c.l.Unlock
}

// ----------------------------------- 状态机 -----------------------------------

type kmStatusI interface {
	record()              // 记录
	playback(name string) // 回放
	pause()               // 暂停
	stop()                // 停止
	save(name string)     // 存储

	current() kmStatusI
}

var chooseBox map[enum.Status]kmStatusI

func init() {
	var base = &BaseT{
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

// --------------------------------- 基础状态 ---------------------------------

type baseStatusT struct {
	base *BaseT
}

func (s *baseStatusT) record() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorFreeToFree],
	})
}
func (s *baseStatusT) playback(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorFreeToFree],
	})
}
func (s *baseStatusT) pause() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorFreeToFree],
	})
}
func (s *baseStatusT) stop() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorFreeToFree],
	})
}
func (s *baseStatusT) save(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorFreeToFree],
	})
}

func (s *baseStatusT) current() kmStatusI {
	return chooseBox[s.base.Status]
}
