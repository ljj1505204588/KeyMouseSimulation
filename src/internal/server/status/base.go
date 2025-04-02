package status

import (
	"KeyMouseSimulation/internal/server/svcComponent"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"fmt"
	"sync"
	"time"
)

type KmStatusI interface {
	Record()              // 记录
	Playback(name string) // 回放
	Pause()               // 暂停
	Stop() (save bool)    // 停止
	Save(name string)     // 存储

	Status() enum.Status
	StatusShow(status enum.Status) string
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
	tryPublishErr(eventCenter.Event.Publish(topic.ServerStatus, &topic.ServerStatusChangeData{
		Status: e,
	}))
}

func (k *kmStatusT) syncServerStatus() {
	defer func() { go k.syncServerStatus() }()

	for range time.NewTicker(1 * time.Second).C {
		tryPublishErr(eventCenter.Event.Publish(topic.ServerStatus, &topic.ServerStatusChangeData{
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
	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.StatusShow(s.getStatus()), language.RecordStr.ToString(), language.ErrorStatusChangeError.ToString()),
	})
}

func (s *baseStatusT) Playback(name string) {
	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.StatusShow(s.getStatus()), language.PlayBackStr.ToString(), language.ErrorStatusChangeError.ToString()),
	})
}

func (s *baseStatusT) Pause() {
	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.StatusShow(s.getStatus()), language.PauseStr.ToString(), language.ErrorStatusChangeError.ToString()),
	})
}

func (s *baseStatusT) Stop() (save bool) {
	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.StatusShow(s.getStatus()), language.StopStr.ToString(), language.ErrorStatusChangeError.ToString()),
	})
	return false
}

func (s *baseStatusT) Save(name string) {
	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: fmt.Sprintf("[%s]->[%s] %s",
			s.StatusShow(s.getStatus()), language.StopStr.ToString(), language.ErrorStatusChangeError.ToString()),
	})
}

func (s *baseStatusT) StatusShow(status enum.Status) string {
	switch status {
	case enum.Free:
		return language.ControlTypeFreeStr.ToString()
	case enum.Recording:
		return language.ControlTypeRecordingStr.ToString()
	case enum.RecordPause:
		return language.ControlTypeRecordPauseStr.ToString()
	case enum.Playback:
		return language.ControlTypePlaybackStr.ToString()
	case enum.PlaybackPause:
		return language.ControlTypePlaybackPauseStr.ToString()
	}
	return ""
}

func tryPublishErr(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{ErrInfo: err.Error()})
	}
}
