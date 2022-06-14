package paramTool

import (
	"KeyMouseSimulation/common/commonTool"
	"os"
	"strconv"
	"strings"
)

var chainMap = map[string]chainIntI{}
var param paramT

type chainIntI interface {
	defaultParamInit(*paramT)
	xmlParamInit(*paramT, paramT)
	envParamInit(*paramT)
}

func chainInit() {
	chainMap = make(map[string]chainIntI)
	chainMap = map[string]chainIntI{
		"ServerParam":  &serverParamT{},
		"EsParam":      &esParamT{},
		"LogParam":     &logParamT{},
		"HotWordParam": &hotWordParamT{},
	}
}

type paramT struct {
	ServerParam  serverParamT
	EsParam      esParamT
	LogParam     logParamT
	HotWordParam hotWordParamT
}

func GetParam() *paramT {
	return &param
}

/*------------------------------------  服务参数  ----------------------------------------------*/
type serverParamT struct {
	ServerPort int      `xml:"ServerPort"`
	Ip         []string `xml:"Ip"`
	OpenPProf  bool     `xml:"OpenPProf"`
	PProfPort  int      `xml:"PProfPort"`
}

func (*serverParamT) defaultParamInit(param *paramT) {
	param.ServerParam.ServerPort = 21002
	param.ServerParam.Ip = []string{
		"111",
		"222",
	}
	param.ServerParam.OpenPProf = false
	param.ServerParam.PProfPort = 6061
}
func (*serverParamT) xmlParamInit(param *paramT, xmlParam paramT) {
	if xmlParam.ServerParam.ServerPort != 0 {
		param.ServerParam.ServerPort = xmlParam.ServerParam.ServerPort
	}
	if xmlParam.ServerParam.OpenPProf == true {
		param.ServerParam.OpenPProf = true
	}
	if xmlParam.ServerParam.PProfPort != 0 {
		param.ServerParam.PProfPort = xmlParam.ServerParam.PProfPort
	}
}
func (*serverParamT) envParamInit(param *paramT) {
	serverPortStr := os.Getenv("ServerParam.ServerPort")
	serverPort, err := strconv.Atoi(serverPortStr)
	if err == nil && serverPort != 0 {
		param.ServerParam.ServerPort = serverPort
	}
	isOpenPProf := os.Getenv("ServerParam.OpenPProf")
	if strings.ToLower(isOpenPProf) == "true" {
		param.ServerParam.OpenPProf = true
	}
	PProfPortStr := os.Getenv("ServerParam.PProfPort")
	PProfPort, err := strconv.Atoi(PProfPortStr)
	if err != nil {
		param.ServerParam.PProfPort = PProfPort
	}
}

/*------------------------------------  日志参数  ------------------------------------------------*/
type logParamT struct {
	LogLevel      string
	LogPath       string
	ExtendLogger  []string
	LogKeepDay    int
	LogRenameSize int
}

func (*logParamT) defaultParamInit(param *paramT) {
	param.LogParam.LogLevel = "INFO"
	param.LogParam.LogPath = commonTool.GetSysPthSep() + "logs"
	param.LogParam.LogKeepDay = 7
	param.LogParam.LogRenameSize = 500
}
func (*logParamT) xmlParamInit(param *paramT, xmlParam paramT) {
	if xmlParam.LogParam.LogLevel != "" {
		param.LogParam.LogLevel = xmlParam.LogParam.LogLevel
	}
	if xmlParam.LogParam.LogPath != "" {
		param.LogParam.LogPath = xmlParam.LogParam.LogPath
	}
	if len(xmlParam.LogParam.ExtendLogger) != 0 {
		param.LogParam.ExtendLogger = xmlParam.LogParam.ExtendLogger
	}
	if xmlParam.LogParam.LogKeepDay != 0 {
		param.LogParam.LogKeepDay = xmlParam.LogParam.LogKeepDay
	}
	if xmlParam.LogParam.LogRenameSize != 0 {
		param.LogParam.LogRenameSize = xmlParam.LogParam.LogRenameSize
	}
}
func (*logParamT) envParamInit(Param *paramT) {
	envLogLevel := os.Getenv("LogParam.LogLevel")
	if envLogLevel != "" {
		Param.LogParam.LogLevel = envLogLevel
	}
	envLogPath := os.Getenv("LogParam.LogPath")
	if envLogPath != "" {
		Param.LogParam.LogPath = envLogPath
	}
	ExtendLogger := os.Getenv("LogParam.ExtendLogger")
	if ExtendLogger != "" {
		loggers := strings.Split(ExtendLogger, ",")
		Param.LogParam.ExtendLogger = loggers
	}
	LogKeepDay := os.Getenv("LogParam.LogKeepDay")
	if LogKeepDay != "" {
		LogKeepDayInt, err := strconv.Atoi(LogKeepDay)
		if err != nil {
			Param.LogParam.LogKeepDay = LogKeepDayInt
		}
	}
	LogRenameSize := os.Getenv("LogParam.LogRenameSize")
	if LogRenameSize != "" {
		LogRenameSizeInt, err := strconv.Atoi(LogRenameSize)
		if err != nil {
			Param.LogParam.LogRenameSize = LogRenameSizeInt
		}
	}
}

/*------------------------------------  Es参数  --------------------------------------------*/
type esParamT struct {
	EsPort     []int  `xml:"EsPort"`
	UserName   string `xml:"UserName"`
	PassWord   string `xml:"PassWord"`
	RecordBody bool   `xml:"RecordBody"`
}

func (*esParamT) defaultParamInit(param *paramT) {
	param.EsParam.EsPort = []int{20003, 21003}
	param.EsParam.UserName = "elastic"
	param.EsParam.PassWord = "jjelastic"
	param.EsParam.RecordBody = false
}
func (*esParamT) xmlParamInit(param *paramT, xmlParam paramT) {
	if len(xmlParam.EsParam.EsPort) != 0 && xmlParam.EsParam.EsPort[0] != 0 {
		param.EsParam.EsPort = xmlParam.EsParam.EsPort
	}
	if xmlParam.EsParam.UserName != "" {
		param.EsParam.UserName = xmlParam.EsParam.UserName
	}
	if xmlParam.EsParam.PassWord != "" {
		param.EsParam.PassWord = xmlParam.EsParam.PassWord
	}
	if xmlParam.EsParam.RecordBody {
		param.EsParam.RecordBody = true
	}
}
func (*esParamT) envParamInit(param *paramT) {
	EsPortStrL := strings.Split(os.Getenv("EsParam.EsPort"), ",")
	var EsPortL = make([]int, 0)
	for _, v := range EsPortStrL {
		Port, err := strconv.Atoi(v)
		if err != nil {
			EsPortL = append(EsPortL, Port)
		}
	}
	if len(EsPortL) != 0 && EsPortL[0] != 0 {
		param.EsParam.EsPort = EsPortL
	}
	if os.Getenv("EsParam.UserName") != "" {
		param.EsParam.UserName = os.Getenv("EsParam.PassWord")
	}
	if os.Getenv("EsParam.UserName") != "" {
		param.EsParam.PassWord = os.Getenv("EsParam.PassWord")
	}
	recordBody := os.Getenv("HotWordParam.RecordBody")
	if strings.ToLower(recordBody) == "true" {
		param.EsParam.RecordBody = true
	}
}

/*------------------------------------  hotWord参数  --------------------------------------------*/
type hotWordParamT struct {
	IsRestart bool
	TimeGap   int
}

func (*hotWordParamT) defaultParamInit(param *paramT) {
	param.HotWordParam.IsRestart = false
	param.HotWordParam.TimeGap = 15
}
func (*hotWordParamT) xmlParamInit(param *paramT, xmlParam paramT) {
	if xmlParam.HotWordParam.IsRestart {
		param.HotWordParam.IsRestart = true
	}
	if xmlParam.HotWordParam.TimeGap != 0 {
		param.HotWordParam.TimeGap = xmlParam.HotWordParam.TimeGap
	}
}
func (*hotWordParamT) envParamInit(param *paramT) {
	isRestart := os.Getenv("HotWordParam.IsRestart")
	if strings.ToLower(isRestart) == "true" {
		param.HotWordParam.IsRestart = true
	}
	timeGapStr := os.Getenv("HotWordParam.TimeGap")
	timeGap, err := strconv.Atoi(timeGapStr)
	if err == nil && timeGap != 0 {
		param.HotWordParam.TimeGap = timeGap
	}
}
