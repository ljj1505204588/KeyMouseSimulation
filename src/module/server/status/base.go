package status

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module/language"
	recordAndPlayBack "KeyMouseSimulation/module/server/BaseComponent"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
)

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
		enum.Free:          &freeStatusT{name: enum.Free, baseStatusT: baseStatusT{base: base}},
		enum.Recording:     &recordingStatusT{name: enum.Recording, baseStatusT: baseStatusT{base: base}},
		enum.RecordPause:   &recordPauseStatusT{name: enum.RecordPause, baseStatusT: baseStatusT{base: base}},
		enum.Playback:      &playbackStatusT{name: enum.Playback, baseStatusT: baseStatusT{base: base}},
		enum.PlaybackPause: &playbackPauseStatusT{name: enum.PlaybackPause, baseStatusT: baseStatusT{base: base}},
	}
}

// --------------------------------- 状态机 ---------------------------------

type BaseT struct {
	Status   enum.Status
	PlayBack recordAndPlayBack.PlayBackServerI
	Record   recordAndPlayBack.RecordServerI
}

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
