package status

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
)

// "RecordPause" 状态机
type recordPauseStatusT struct {
	name enum.Status
	baseStatusT
}

func (s *recordPauseStatusT) record() {
	s.base.Record.Start()
	s.base.Status = enum.Recording
}
func (s *recordPauseStatusT) stop() {
	s.base.Record.Stop()
	s.base.Status = enum.Free
	_ = eventCenter.Event.Publish(events.RecordFinish, events.RecordFinishData{})
}
func (s *recordPauseStatusT) status() enum.Status {
	return enum.RecordPause
}
