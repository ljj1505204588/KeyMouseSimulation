package ui

import "KeyMouseSimulation/module/server"

const (
	UI_TYPE_FREE           = int(server.CONTROL_TYPE_FREE)
	UI_TYPE_RECORDING      = int(server.CONTROL_TYPE_RECORDING)
	UI_TYPE_RECORD_PAUSE   = int(server.CONTROL_TYPE_RECORD_PAUSE)
	UI_TYPE_PLAYBACK       = int(server.CONTROL_TYPE_PLAYBACK)
	UI_TYPE_PLAYBACK_PAUSE = int(server.CONTROL_TYPE_PLAYBACK_PAUSE)
)

