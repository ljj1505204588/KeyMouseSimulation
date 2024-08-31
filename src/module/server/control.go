package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/module/server/status"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
)

func init() {
	eventCenter.Event.Register(events.ButtonClick, buttonClickHandler)
}

var control = status.NewKmStatusI()

func buttonClickHandler(data interface{}) (err error) {
	var d = data.(events.ButtonClickData)
	switch d.Button {
	case enum.RecordButton:
		control.Record()
	case enum.PlaybackButton:
		control.Playback(d.Name)
	case enum.PauseButton:
		control.Pause()
	case enum.StopButton:
		control.Stop()
	case enum.SaveFileButton:
		control.Save(d.Name)
	}
	return
}
