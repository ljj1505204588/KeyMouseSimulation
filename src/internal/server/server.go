package server

import (
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
)

var Svc SvcI = &svcT{}

type SvcI interface {
	StatusShow(status enum.Status) string
	Record()           // 记录
	PlayBack()         // 回放
	Pause()            // 暂停
	Stop() (save bool) // 停止
	Save(name string)  // 存储文件
}

type svcT struct {
}

func (s *svcT) StatusShow(status enum.Status) string {
	switch status {
	case enum.Free:
		return language.ControlTypeFreeStr.ToString()
	case enum.Recording:
		return language.ControlTypeRecordingStr.ToString()
	case enum.RecordPause:
		return language.ControlTypeRecordPauseStr.ToString()
	case enum.Playback:
		return language.ControlTypePlaybackStr.ToString()
	case enum.PlaybackPause:
		return language.ControlTypePlaybackPauseStr.ToString()
	}
	return ""
}

func (s *svcT) Record() {

}
func (s *svcT) PlayBack() {

}
func (s *svcT) Pause() {

}

// Save 存储文件
func (s *svcT) Save(name string) {

}
func (s *svcT) Stop() (save bool) {
	return
}
