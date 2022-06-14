package logTool

import (
	"KeyMouseSimulation/common/paramTool"
	"os"
	"strconv"
)

const LOG_TOOL_LOG_HEAD = "日志工具："

var ajLogController struct {
	serverLog *ajLoggerT
	accessLog *ajLoggerT
	extendLog map[string]*ajLoggerT

	logLevel int
}

func init() {
	createPidFile()

	var err error
	//所有日志输出初始化、自动管理
	ajLogController.serverLog, err = NewAJLogger(SERVER_LOG_NAME)
	if err != nil {
		panic(LOG_TOOL_LOG_HEAD + "创建日志Logger出错." + err.Error())
	}
	ajLogController.accessLog, err = NewAJLogger(ACCESS_LOG_NAME)
	if err != nil {
		panic(LOG_TOOL_LOG_HEAD + "创建日志Logger出错." + err.Error())
	}

	go AutoManagementAJLogger(ajLogController.accessLog)
	go AutoManagementAJLogger(ajLogController.serverLog)

	for _, v := range paramTool.GetParam().LogParam.ExtendLogger {
		ajLogController.extendLog[v], err = NewAJLogger(v)
		if err != nil {
			panic(LOG_TOOL_LOG_HEAD + "创建日志Logger出错." + err.Error())
		}
		go AutoManagementAJLogger(ajLogController.extendLog[v])
	}

	logLevelMapInit()
	logLevelInt, ok := logLevelMap[paramTool.GetParam().LogParam.LogLevel]
	if !ok {
		panic(LOG_TOOL_LOG_HEAD + "配置中默认日志等级不在可选范围 【DEBUG INFO WARNING ERROR FATAL】内.")
	}
	ajLogController.logLevel = logLevelInt
	InfoAJ(LOG_TOOL_LOG_HEAD + "日志组件，加载成功！")

	//初始化日志等级
	_, ok = logLevelMap[paramTool.GetParam().LogParam.LogLevel]
	if !ok {
		panic("配置中日志等级错误.当前配置：" + paramTool.GetParam().LogParam.LogLevel)
	}
}

/*----------------------日志等级输出---------------------------*/

func DebugAJ(msg string, outPutLogger ...string) {
	if ajLogController.logLevel > lEVEL_DEBUG_INT {
		return
	}
	reCordLog(LEVEL_DEBUG_STR+" "+msg, outPutLogger...)
	return
}
func InfoAJ(msg string, outPutLogger ...string) {
	if ajLogController.logLevel > lEVEL_INFO_INT {
		return
	}
	reCordLog(LEVEL_INFO_STR+" "+msg, outPutLogger...)
	return
}
func ErrorAJ(err error, outPutLogger ...string) {
	if ajLogController.logLevel > lEVEL_ERROR_INT || err == nil {
		return
	}
	reCordLog(LEVEL_ERROR_STR+" "+err.Error(), outPutLogger...)
	return
}
func WarningAJ(msg string, outPutLogger ...string) {
	if ajLogController.logLevel > lEVEL_WARNING_INT {
		return
	}
	reCordLog(LEVEL_WARNING_STR+" "+msg, outPutLogger...)
	return
}
func FatalAJ(msg string, outPutLogger ...string) {
	reCordLog(FATAL_LEVEL_STR+" "+msg, outPutLogger...)
	os.Exit(1)
}

/*-----------------------------------------------------------*/

func SetLogLevel(logLevel string) error {
	logLevelInt, ok := logLevelMap[logLevel]
	if !ok {
		return NewLogError(NIL_ERROR_CODE, LOG_TOOL_LOG_HEAD+"配置中默认日志等级不在可选范围 【DEBUG INFO WARNING ERROR FATAL】内.", nil)
	}
	ajLogController.logLevel = logLevelInt
	return nil
}
func reCordLog(msg string, outPutLogger ...string) {
	if len(outPutLogger) == 0 {
		_ = ajLogController.serverLog.Logger.Output(3, msg)
		return
	}
	for _, v := range outPutLogger {
		if v == ACCESS_LOG_NAME {
			_ = ajLogController.accessLog.Logger.Output(5, msg)
		} else {
			AjLogger, ok := ajLogController.extendLog[v]
			if ok {
				_ = AjLogger.Logger.Output(3, msg)
			}
		}
	}
}
func createPidFile() {
	pid := os.Getpid()
	file, _ := os.OpenFile("pid.txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0762)
	_, _ = file.Write([]byte(strconv.Itoa(pid)))
	_ = file.Close()
}
