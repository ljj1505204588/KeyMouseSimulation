package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/commonTool"
	"KeyMouseSimulation/common/share/enum"
	"KeyMouseSimulation/common/share/events"
	"KeyMouseSimulation/module/server/status"
	"sync"
)

func init() {
	eventCenter.Event.Register(events.ButtonClick, server.buttonClickHandler)
	eventCenter.Event.Register(events.PlaybackFinish, server.playbackHandler)
}

var server = &serverT{
	control: status.NewKmStatusI(),
}

type serverT struct {
	lock    sync.Mutex
	control status.KmStatusI
}

func (s *serverT) playbackHandler(dataI interface{}) (err error) {
	defer commonTool.LockSelf(&s.lock)()

	s.control.Stop()
	return
}
func (s *serverT) buttonClickHandler(dataI interface{}) (err error) {
	if !s.lock.TryLock() {
		return
	}
	defer s.lock.Unlock()

	var data = dataI.(events.ButtonClickData)
	switch data.Button {
	case enum.RecordButton:
		s.control.Record()
	case enum.PlaybackButton:
		s.control.Playback(data.Name)
	case enum.PauseButton:
		s.control.Pause()
	case enum.StopButton:
		s.control.Stop()
	}
	return
}
