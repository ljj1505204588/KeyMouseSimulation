package status

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/share/enum"
	"KeyMouseSimulation/common/share/events"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/module/server/svcComponent"
	"fmt"
	"sync"
	"time"
)

type KmStatusI interface {
	Record()              // 记录
	Playback(name string) // 回放
	Pause()               // 暂停
	Stop()                // 停止

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
		getStatus: k.getStatus,
		playBack:  svcComponent.GetPlaybackServer(),
		record:    svcComponent.GetRecordServer(),
	}
	k.statusBox = map[enum.Status]KmStatusI{
		enum.Free:          &freeStatusT{name: enum.Free, baseStatusT: baseStatus},
		enum.Recording:     &recordingStatusT{name: enum.Recording, baseStatusT: baseStatus},
		enum.RecordPause:   &recordPauseStatusT{name: enum.RecordPause, baseStatusT: baseStatus},
		enum.Playback:      &playbackStatusT{name: enum.Playback, baseStatusT: baseStatus},
		enum.PlaybackPause: &playbackPauseStatusT{name: enum.PlaybackPause, baseStatusT: baseStatus},
	}
	k.setStatus(enum.Free)
	go k.syncServerStatus()
}

func (k *kmStatusT) getStatus() enum.Status {
	return k.KmStatusI.Status()
}

func (k *kmStatusT) setStatus(e enum.Status) {
	defer k.lockSelf()()
	k.KmStatusI = k.statusBox[e]
	tryPublishErr(eventCenter.Event.Publish(events.ServerStatus, events.ServerStatusChangeData{
		Status: e,
	}))
}

func (k *kmStatusT) syncServerStatus() {
	defer func() { go k.syncServerStatus() }()

	for range time.NewTicker(1 * time.Second).C {
		tryPublishErr(eventCenter.Event.Publish(events.ServerStatus, events.ServerStatusChangeData{
			Status: k.KmStatusI.Status(),
		}))
	}
}

func (k *kmStatusT) lockSelf() func() {
	k.lock.Lock()
	return k.lock.Unlock
}

// --------------------------------- 状态机 ---------------------------------

type baseStatusT struct {
	setStatus func(enum.Status)
	getStatus func() enum.Status
	playBack  svcComponent.PlayBackServerI
	record    svcComponent.RecordServerI
}

func (s *baseStatusT) Record() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.getStatus().Language(), language.Center.Get(language.RecordStr), language.Center.Get(language.ErrorStatusChangeError)),
	})
}

func (s *baseStatusT) Playback(name string) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.getStatus().Language(), language.Center.Get(language.PlaybackStr), language.Center.Get(language.ErrorStatusChangeError)),
	})
}

func (s *baseStatusT) Pause() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.getStatus().Language(), language.Center.Get(language.PauseStr), language.Center.Get(language.ErrorStatusChangeError)),
	})
}

func (s *baseStatusT) Stop() {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.getStatus().Language(), language.Center.Get(language.StopStr), language.Center.Get(language.ErrorStatusChangeError)),
	})
}

func tryPublishErr(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
	}
}
