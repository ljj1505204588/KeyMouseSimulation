package logTool

var logLevelMap map[string]int

/*------------------- 日志等级 ----------------------*/
func logLevelMapInit() {
	logLevelMap = make(map[string]int)
	logLevelMap = map[string]int{
		LEVEL_DEBUG_STR:   lEVEL_DEBUG_INT,
		LEVEL_INFO_STR:    lEVEL_INFO_INT,
		LEVEL_WARNING_STR: lEVEL_WARNING_INT,
		LEVEL_ERROR_STR:   lEVEL_ERROR_INT,
		FATAL_LEVEL_STR:   lEVEL_FATAL_INT,
	}
}

const (
	lEVEL_DEBUG_INT = 1 << iota
	lEVEL_INFO_INT
	lEVEL_WARNING_INT
	lEVEL_ERROR_INT
	lEVEL_FATAL_INT
)
const (
	LEVEL_DEBUG_STR   = "DEBUG"
	LEVEL_INFO_STR    = "INFO"
	LEVEL_WARNING_STR = "WARNING"
	LEVEL_ERROR_STR   = "ERROR"
	FATAL_LEVEL_STR   = "FATAL"
)

/*-------------------  其他  ----------------------*/
const (
	LOG_NAME_TAIL   = ".log"
	SERVER_LOG_NAME = "server"
	ACCESS_LOG_NAME = "access"
)
