package status

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
)

// "Recording" 状态机
type recordingStatusT struct {
	name enum.Status
	baseStatusT
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
