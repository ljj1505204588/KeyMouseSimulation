package status

import "KeyMouseSimulation/share/enum"

// "Free" 状态机
type freeStatusT struct {
	name enum.Status
	baseStatusT
}

func (s *freeStatusT) record() {
	s.base.Record.Start()
	s.base.Status = enum.Recording
}
func (s *freeStatusT) playback(name string) {
	s.base.PlayBack.Start(name)
	s.base.Status = enum.Playback
}
func (s *freeStatusT) save(name string) {
	s.base.Record.Save(name)
}
func (s *freeStatusT) status() enum.Status {
	return enum.Free
}
