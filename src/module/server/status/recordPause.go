package status

import (
	"KeyMouseSimulation/share/enum"
)

// "RecordPause" 状态机
type recordPauseStatusT struct {
	name enum.Status
	*baseStatusT
}

func (s *recordPauseStatusT) Record() {
	s.record.Start()
	s.setStatus(enum.Recording)
}
func (s *recordPauseStatusT) Stop() {
	s.record.Stop()
	s.setStatus(enum.Free)
}
func (s *recordPauseStatusT) Status() enum.Status {
	return enum.RecordPause
}
