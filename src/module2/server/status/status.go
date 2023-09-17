package status

import (
	"KeyMouseSimulation/module2/server/BaseComponent"
	"KeyMouseSimulation/share/enum"
)

type KMStatusI interface {
	Record()   // 记录
	Playback() // 回放
	Pause()    // 暂停
	Stop()     // 停止
	Save()     // 存储
}

var chooseBox map[enum.Status]KMStatusI

func InitKMStatus(base *recordAndPlayBack.BaseT) KMStatusI {
	chooseBox = map[enum.Status]KMStatusI{
		enum.Free:          &FreeStatusT{&BaseStatusT{base: base}},
		enum.Recording:     &RecordingStatusT{&BaseStatusT{base: base}},
		enum.RecordPause:   &RecordPauseStatusT{&BaseStatusT{base: base}},
		enum.Playback:      &PlaybackStatusT{&BaseStatusT{base: base}},
		enum.PlaybackPause: &PlaybackPauseStatusT{&BaseStatusT{base: base}},
	}

	return chooseBox[enum.Free]
}

// --------------------------------- 状态机 ---------------------------------

// FreeStatusT "Free" 状态机
type FreeStatusT struct {
	*BaseStatusT
}

func (s *FreeStatusT) Record() {
	s.base.Record.Start()
}
func (s *FreeStatusT) Playback() {
	s.base.PlayBack.Start()
}
func (s *FreeStatusT) Save() {
	s.base.Record.Save()
}

// RecordingStatusT "Recording" 状态机
type RecordingStatusT struct {
	*BaseStatusT
}

func (s *BaseStatusT) Record()   {}
func (s *BaseStatusT) Playback() {}
func (s *BaseStatusT) Pause()    {}
func (s *BaseStatusT) Stop()     {}
func (s *BaseStatusT) Save()     {}

// RecordPauseStatusT "RecordPause" 状态机
type RecordPauseStatusT struct {
	*BaseStatusT
}

func (s *BaseStatusT) Record()   {}
func (s *BaseStatusT) Playback() {}
func (s *BaseStatusT) Pause()    {}
func (s *BaseStatusT) Stop()     {}
func (s *BaseStatusT) Save()     {}

// PlaybackStatusT "Playback" 状态机
type PlaybackStatusT struct {
	*BaseStatusT
}

func (s *BaseStatusT) Record()   {}
func (s *BaseStatusT) Playback() {}
func (s *BaseStatusT) Pause()    {}
func (s *BaseStatusT) Stop()     {}
func (s *BaseStatusT) Save()     {}

// PlaybackPauseStatusT "PlaybackPause" 状态机
type PlaybackPauseStatusT struct {
	*BaseStatusT
}

func (s *BaseStatusT) Record()   {}
func (s *BaseStatusT) Playback() {}
func (s *BaseStatusT) Pause()    {}
func (s *BaseStatusT) Stop()     {}
func (s *BaseStatusT) Save()     {}

// --------------------------------- 基础状态 ---------------------------------

type BaseStatusT struct {
	base *recordAndPlayBack.BaseT
}

func (s *BaseStatusT) Record()   {}
func (s *BaseStatusT) Playback() {}
func (s *BaseStatusT) Pause()    {}
func (s *BaseStatusT) Stop()     {}
func (s *BaseStatusT) Save()     {}
