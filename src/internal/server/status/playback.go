package status

import "KeyMouseSimulation/share/enum"

// "Playback" 状态机
type playbackStatusT struct {
	name enum.Status
	*baseStatusT
}

func (s *playbackStatusT) Pause() {
	s.playBack.Pause()
	s.setStatus(enum.PlaybackPause)
}
func (s *playbackStatusT) Stop() (save bool) {
	s.playBack.Stop()
	s.setStatus(enum.Free)
	return false
}
func (s *playbackStatusT) Status() enum.Status {
	return enum.Playback
}
