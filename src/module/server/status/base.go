package status

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module/language"
	recordAndPlayBack "KeyMouseSimulation/module/server/BaseComponent"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"sync"
)

type KmStatusI interface {
	Record()              // 记录
	Playback(name string) // 回放
	Pause()               // 暂停
	Stop()                // 停止
	Save(name string)     // 存储

	Status() enum.Status
}

func NewKmStatusI() KmStatusI {
	var kmT = kmStatusT{}
	kmT.init()
	return &kmT
}

type kmStatusT struct {
	KmStatusI
	statusBox map[enum.Status]KmStatusI

	lock sync.Mutex
}

func (k *kmStatusT) init() {

	var baseStatus = &baseStatusT{
		setStatus: k.setStatus,
		playBack:  recordAndPlayBack.GetPlaybackServer(),
		record:    recordAndPlayBack.GetRecordServer(),
	}
	k.statusBox = map[enum.Status]KmStatusI{
		enum.Free:          &freeStatusT{name: enum.Free, baseStatusT: baseStatus},
		enum.Recording:     &recordingStatusT{name: enum.Recording, baseStatusT: baseStatus},
		enum.RecordPause:   &recordPauseStatusT{name: enum.RecordPause, baseStatusT: baseStatus},
		enum.Playback:      &playbackStatusT{name: enum.Playback, baseStatusT: baseStatus},
		enum.PlaybackPause: &playbackPauseStatusT{name: enum.PlaybackPause, baseStatusT: baseStatus},
	}
	k.KmStatusI = k.statusBox[enum.Free]
}

func (k *kmStatusT) setStatus(e enum.Status) {
	defer k.lockSelf()()
	k.KmStatusI = k.statusBox[e]
}

func (k *kmStatusT) lockSelf() func() {
	k.lock.Lock()
	return k.lock.Unlock
}

// --------------------------------- 状态机 ---------------------------------

type baseStatusT struct {
	setStatus func(enum.Status)
	playBack  recordAndPlayBack.PlayBackServerI
	record    recordAndPlayBack.RecordServerI
}

func (s *baseStatusT) Record() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) Playback(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) Pause() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) Stop() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
func (s *baseStatusT) Save(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: language.CurrentUse[language.ErrorStatusChangeError],
	})
}
