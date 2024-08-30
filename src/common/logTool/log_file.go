package logTool

import (
	"KeyMouseSimulation/common/commonTool"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ajLoggerT struct {
	Logger  *log.Logger
	LogFile *os.File
}

func (T *ajLoggerT) Write(msg []byte) (int, error) {
	_, err := T.LogFile.Write(msg)
	if err != nil {
		return -1, err
	}
	fmt.Println(string(msg))
	return 0, nil
}

func NewAJLogger(name string) (*ajLoggerT, error) {
	wd, _ := os.Getwd()

	var logPath string
	if LogParam.LogPath != "" {
		logPath = wd + commonTool.GetSysPthSep() + LogParam.LogPath
	} else {
		logPath = wd + commonTool.GetSysPthSep() + "logs"
	}
	err := os.MkdirAll(logPath, 0762)
	if err != nil {
		return nil, err
	}

	filePath := logPath + commonTool.GetSysPthSep() + name + LOG_NAME_TAIL
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}

	AJLog := &ajLoggerT{
		LogFile: logFile,
	}

	logger := log.New(AJLog, "", log.Ldate|log.Ltime|log.Lshortfile)
	AJLog.Logger = logger

	return AJLog, nil
}

func AutoManagementAJLogger(AJLogger *ajLoggerT) {
	if AJLogger == nil {
		return
	}
	<-getZeroTimer().C
	timer := time.NewTicker(24 * time.Hour)
	for {
		timer.Reset(24 * time.Hour)
		//给旧文件修改名称，绑定新文件
		dir, name := filepath.Split(AJLogger.LogFile.Name())
		nameWithoutTail := name[:len(name)-len(LOG_NAME_TAIL)]

		//创建临时日志输出，关闭原有日志
		tmpFilePath := dir + commonTool.GetSysPthSep() + nameWithoutTail + "-tmp" + LOG_NAME_TAIL
		tmpFile, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			panic(err.Error())
		}

		AJLogger.Logger.SetOutput(tmpFile)
		_ = AJLogger.LogFile.Close()

		//修改日志名称
		reNamePath := dir + commonTool.GetSysPthSep() + nameWithoutTail + time.Now().Format("-2006-01-02") + LOG_NAME_TAIL
		_ = os.Rename(AJLogger.LogFile.Name(), reNamePath)

		//绑定新的输出日志
		file, _ := os.OpenFile(dir+commonTool.GetSysPthSep()+name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		AJLogger.Logger.SetOutput(file)
		_ = tmpFile.Close()
		AJLogger.LogFile = file

		//读取临时文件中的日志存储到前一天日志中。删除临时日志
		tmpInfo, _ := os.ReadFile(tmpFilePath)
		reNameFile, _ := os.OpenFile(reNamePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		_, _ = reNameFile.Write(tmpInfo)
		_ = reNameFile.Close()
		_ = os.Remove(tmpFilePath)

		//删除过期无效数据
		files, _ := ioutil.ReadDir(dir)
		for _, v := range files {
			data, err := time.Parse(nameWithoutTail+"-2006-01-02"+LOG_NAME_TAIL, v.Name())
			if err == nil && time.Now().Unix()-data.Unix() > int64(3600*24*LogParam.LogKeepDay) {
				_ = os.Remove(dir + commonTool.GetSysPthSep() + v.Name())
			}
		}

		<-timer.C
	}
}
func getZeroTimer() *time.Timer {
	now := time.Now()
	y, m, d := now.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	return time.NewTimer(time.Duration(nextDay.UnixNano()-now.UnixNano()) - 1*time.Minute)
}
