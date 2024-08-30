package status

import "KeyMouseSimulation/share/enum"

// "Playback" 状态机
type playbackStatusT struct {
	name enum.Status
	baseStatusT
}

func (s *playbackStatusT) pause() {
	s.base.PlayBack.Pause()
	s.base.Status = enum.PlaybackPause
}
func (s *playbackStatusT) stop() {
	s.base.PlayBack.Stop()
	s.base.Status = enum.Free
}
func (s *playbackStatusT) status() enum.Status {
	return enum.Playback
}
