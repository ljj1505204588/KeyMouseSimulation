package server

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/internal/server/status"
	eventCenter "KeyMouseSimulation/pkg/event"
	rp_file "KeyMouseSimulation/pkg/file"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"sync"
)

func init() {
	var svc = Svc.(*serverT)
	eventCenter.Event.Register(topic.PlaybackFinish, svc.playbackHandler)
}

type serverT struct {
	lock    sync.Mutex
	control status.KmStatusI
}

func (s *serverT) playbackHandler(dataI interface{}) (err error) {
	defer common.LockSelf(&s.lock)()

	s.control.Stop()
	return
}

func (s *serverT) StatusShow(status enum.Status) string {
	return s.control.StatusShow(status)
}

func (s *serverT) Record() {
	defer common.LockSelf(&s.lock)()
	s.control.Record()
}
func (s *serverT) PlayBack() {
	defer common.LockSelf(&s.lock)()

	var current = rp_file.FileControl.Current()
	s.control.Playback(current)
}
func (s *serverT) Pause() {
	defer common.LockSelf(&s.lock)()
	s.control.Pause()
}

// Save 存储文件
func (s *serverT) Save(name string) {
	defer common.LockSelf(&s.lock)()
	s.control.Save(name)
}
func (s *serverT) Stop() (save bool) {
	defer common.LockSelf(&s.lock)()
	return s.control.Stop()
}
