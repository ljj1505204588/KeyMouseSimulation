package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module/server/status"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"sync"
)

func init() {
	eventCenter.Event.Register(events.ButtonClick, server.buttonClickHandler)
}

var server = &serverT{
	control: status.NewKmStatusI(),
}

type serverT struct {
	control status.KmStatusI
	lock    sync.Mutex
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
	case enum.SaveFileButton:
		s.control.Save(data.Name)
	}
	return
}
