package status

import "KeyMouseSimulation/share/enum"

// "PlaybackPause" 状态机
type playbackPauseStatusT struct {
	name enum.Status
	*baseStatusT
}

func (s *playbackPauseStatusT) Playback(name string) {
	s.playBack.Start(name)
	s.setStatus(enum.Playback)
}
func (s *playbackPauseStatusT) Stop() (save bool) {
	s.playBack.Stop()
	s.setStatus(enum.Free)
	return false
}
func (s *playbackPauseStatusT) Status() enum.Status {
	return enum.PlaybackPause
}
