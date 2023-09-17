package paramTool

import (
	"KeyMouseSimulation/common/commonTool"
	"os"
	"strconv"
	"strings"
)

/*
添加参数类型
step1:添加实现`模板接口`的类
step2:添加到paramT类中
step3:添加对应枚举 -- 如果需要默认激活着添加到 defaultActive 中
*默认激活 ： 会回写到配置文件中

//todo 加一个加密以及xml文件中注释的功能·
*/

/*------------------------------------  参数结构体 && 枚举  ----------------------------------------------*/

type paramT struct {
	ServerParam    *serverParamT    `enum:"ServerParam"`
	EsParam        *esParamT        `enum:"EsParam"`
	LogParam       *logParamT       `enum:"LogParam"`
	HotWordParam   *hotWordParamT   `enum:"HotWordParam"`
	MongodbParam   *mongodbParamT   `enum:"MongodbParam"`
	NsqParam       *nsqParamT       `enum:"NsqParam"`
	HeartbeatParam *heartbeatParamT `enum:"HeartbeatParam"`
	RedisParam     *redisParamT     `enum:"RedisParam"`
}

type ParamEnum string

const (
	ServerParamEnum    ParamEnum = "ServerParam"
	EsParamEnum        ParamEnum = "EsParam"
	LogParamEnum       ParamEnum = "LogParam"
	HotWordParamEnum   ParamEnum = "HotWordParam"
	MongodbParamEnum   ParamEnum = "MongodbParam"
	NsqParamEnum       ParamEnum = "NsqParam"
	RedisParam         ParamEnum = "RedisParam"
	HeartbeatParamEnum ParamEnum = "HeartbeatParam"
)

var defaultActive = map[ParamEnum]bool{
	LogParamEnum: true,
}

// 模板
type templateT interface {
	EnvParamInit(*paramT)     // 环境变量配置读取
	DefaultParamInit(*paramT) // 默认配置读取
}

/*------------------------------------  服务参数  ----------------------------------------------*/
type serverParamT struct {
	ServerPort int  `xml:"ServerPort"`
	OpenPProf  bool `xml:"OpenPProf"`
	PProfPort  int  `xml:"PProfPort"`
}

func (*serverParamT) EnvParamInit(param *paramT) {
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
	if err == nil && PProfPort != 0 {
		param.ServerParam.PProfPort = PProfPort
	}
}
func (*serverParamT) DefaultParamInit(param *paramT) {
	if param.ServerParam.ServerPort == 0 {
		param.ServerParam.ServerPort = 21002
	}
	if param.ServerParam.PProfPort == 0 {
		param.ServerParam.PProfPort = 6061
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

func (*logParamT) EnvParamInit(Param *paramT) {
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
func (*logParamT) DefaultParamInit(param *paramT) {
	if param.LogParam.LogLevel == "" {
		param.LogParam.LogLevel = "INFO"
	}
	if param.LogParam.LogPath == "" {
		param.LogParam.LogPath = commonTool.GetSysPthSep() + "logs"
	}
	if param.LogParam.LogKeepDay == 0 {
		param.LogParam.LogKeepDay = 7
	}
	if param.LogParam.LogRenameSize == 0 {
		param.LogParam.LogRenameSize = 500
	}
}

/*------------------------------------  Es参数  --------------------------------------------*/
type esParamT struct {
	EsPort     []int  `xml:"EsPort"`
	UserName   string `xml:"UserName"`
	PassWord   string `xml:"PassWord"`
	RecordBody bool   `xml:"RecordBody"`
}

func (*esParamT) EnvParamInit(param *paramT) {
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
func (*esParamT) DefaultParamInit(param *paramT) {
	if len(param.EsParam.EsPort) == 0 {
		param.EsParam.EsPort = []int{20003, 21003}
	}
	if param.EsParam.UserName == "" {
		param.EsParam.UserName = "elastic"
	}
	if param.EsParam.PassWord == "" {
		param.EsParam.PassWord = "jjelastic"
	}
}

/*------------------------------------  hotWord参数  --------------------------------------------*/
type hotWordParamT struct {
	IsRestart bool
	TimeGap   int
}

func (*hotWordParamT) EnvParamInit(param *paramT) {
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
func (*hotWordParamT) DefaultParamInit(param *paramT) {
	if param.HotWordParam.TimeGap == 0 {
		param.HotWordParam.TimeGap = 15
	}
}

/*------------------------------------  mongodb参数  --------------------------------------------*/
type mongodbParamT struct {
	Uri      string // 连接串
	DataBase string // 数据库
	ReadDay  int    // 读取天数 点触需求，后续删除
}

func (*mongodbParamT) EnvParamInit(param *paramT) {
	uri := os.Getenv("MongodbParam.Uri")
	if uri != "" {
		param.MongodbParam.Uri = uri
	}
	collect := os.Getenv("MongodbParam.DataBase")
	if collect != "" {
		param.MongodbParam.DataBase = collect
	}
}
func (*mongodbParamT) DefaultParamInit(param *paramT) {
	if param.MongodbParam.Uri == "" {
		param.MongodbParam.Uri = ""
	}
	if param.MongodbParam.DataBase == "" {
		param.MongodbParam.DataBase = ""
	}
	if param.MongodbParam.ReadDay == 0 {
		param.MongodbParam.ReadDay = 2
	}
}

/*------------------------------------  nsq参数  --------------------------------------------*/
type nsqParamT struct {
	NsqUri       []string // 连接串
	NsqLookupUrl string   // 管理连接器
}

func (*nsqParamT) EnvParamInit(param *paramT) {
	uriStr := os.Getenv("NsqParam.NsqUrl")
	if uriStr != "" {
		var url = strings.Split(uriStr, ",")
		param.NsqParam.NsqUri = url
	}
	lookupUrl := os.Getenv("NsqParam.NsqLookupUrl")
	if lookupUrl != "" {
		param.NsqParam.NsqLookupUrl = lookupUrl
	}
}
func (*nsqParamT) DefaultParamInit(param *paramT) {
	if len(param.NsqParam.NsqUri) == 0 {
		param.NsqParam.NsqUri = []string{"127.0.0.1:4150"}
	}
	if param.NsqParam.NsqLookupUrl == "" {
		param.NsqParam.NsqLookupUrl = "127.0.0.1:4161"
	}
}

/*------------------------------------  redis  --------------------------------------------*/

type redisParamT struct {
	Cluster  bool     // 是否使用集群
	MulAddr  []string // 地址
	PassWord string   // 密码
}

func (*redisParamT) EnvParamInit(param *paramT) {

}

func (*redisParamT) DefaultParamInit(param *paramT) {

}

/*------------------------------------  心跳参数  --------------------------------------------*/

type heartbeatParamT struct {
	Topic      string
	ServerName string

	GapTime     int // 间隔时间
	OfflineTime int // 离线时长

	ManagerMail string
}

func (*heartbeatParamT) EnvParamInit(param *paramT) {

}

func (*heartbeatParamT) DefaultParamInit(param *paramT) {
	if param.HeartbeatParam.GapTime == 0 {
		param.HeartbeatParam.GapTime = 5
	}
	if param.HeartbeatParam.OfflineTime == 0 {
		param.HeartbeatParam.OfflineTime = 90
	}
}
