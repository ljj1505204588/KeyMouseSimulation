package status

import (
	conf "KeyMouseSimulation/pkg/config"
	"KeyMouseSimulation/share/enum"
)

// "Free" 状态机
type freeStatusT struct {
	name enum.Status
	*baseStatusT
}

func (s *freeStatusT) Record() {
	s.record.Start()
	s.setStatus(enum.Recording)
}
func (s *freeStatusT) Playback(name string) {
	// 重置回放次数
	var playbackTimes = conf.PlaybackTimesConf.GetValue()
	conf.PlaybackRemainTimesConf.SetValue(playbackTimes)

	s.playBack.Start(name)
	s.setStatus(enum.Playback)
}
func (s *freeStatusT) Status() enum.Status {
	return enum.Free
}

func (s *freeStatusT) Save(name string) {
	s.record.Save(name)
}
