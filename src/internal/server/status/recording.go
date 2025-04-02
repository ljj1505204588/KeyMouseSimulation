package status

import "KeyMouseSimulation/share/enum"

// "Recording" 状态机
type recordingStatusT struct {
	name enum.Status
	*baseStatusT
}

func (s *recordingStatusT) Pause() {
	s.record.Pause()
	s.setStatus(enum.RecordPause)
}
func (s *recordingStatusT) Stop() (save bool) {
	s.record.Stop()
	s.setStatus(enum.Free)
	return true
}
func (s *recordingStatusT) Status() enum.Status {
	return enum.Recording
}
