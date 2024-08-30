package logTool

import (
	"KeyMouseSimulation/common/commonTool"
)

var LogParam logParamT

type logParamT struct {
	LogLevel      string
	LogPath       string
	ExtendLogger  []string
	LogKeepDay    int
	LogRenameSize int
}

func (l *logParamT) defaultParam() {
	l.LogLevel = "INFO"
	l.LogPath = commonTool.GetSysPthSep() + "logs"
	l.LogKeepDay = 7
	l.LogRenameSize = 500
}

func (l *logParamT) ParamName() (name string) {
	return "logParam"
}
