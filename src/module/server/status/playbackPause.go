package status

import "KeyMouseSimulation/share/enum"

// "PlaybackPause" 状态机
type playbackPauseStatusT struct {
	name enum.Status
	baseStatusT
}

func (s *playbackPauseStatusT) playback(name string) {
	s.base.PlayBack.Start(name)
	s.base.Status = enum.Playback
}
func (s *playbackPauseStatusT) stop() {
	s.base.PlayBack.Stop()
	s.base.Status = enum.Free
}
func (s *playbackPauseStatusT) status() enum.Status {
	return enum.PlaybackPause
}
