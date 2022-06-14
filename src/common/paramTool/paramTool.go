package paramTool

import (
	"KeyMouseSimulation/common/commonTool"
	"encoding/xml"
	"fmt"
	"os"
)

const PARAM_TOOL_LOG_HEAD = "配置读取工具："

func init() {
	// 打开文件
	fInfo := getFileInfo()

	//读取配置
	var xmlParam paramT
	_ = xml.Unmarshal(fInfo, &xmlParam)
	writeBack := readConfig(xmlParam)

	//复写文件
	f := getFile(os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	defer func() {
		_ = f.Close()
	}()
	_, _ = f.Write(writeBack)

	fmt.Println(PARAM_TOOL_LOG_HEAD + "参数组件，加载成功！")
}

//readConfig 初始化责任链 读取相关配置
func readConfig(xmlParam paramT) []byte {
	chainInit()

	for _, v := range chainMap {
		v.defaultParamInit(&param)
		v.xmlParamInit(&param, xmlParam)
		v.envParamInit(&param)
	}

	writeBack, err := xml.MarshalIndent(param, "", "    ")
	if err != nil {
		panic(PARAM_TOOL_LOG_HEAD + "复写回文件配置编码 Xml 失败.")
	}

	return writeBack
}

func getFile(flag int) *os.File {
	wd, _ := os.Getwd()

	configPath := wd + commonTool.GetSysPthSep() + "config"
	_ = os.Mkdir(configPath, 0764)

	filePath := configPath + commonTool.GetSysPthSep() + "config.xml"

	file, err := os.OpenFile(filePath, flag, 0762)
	if err != nil {
		panic(PARAM_TOOL_LOG_HEAD + "创建文件时错误." + err.Error())
	}

	return file
}
func getFileInfo() []byte {
	wd, _ := os.Getwd()

	configPath := wd + commonTool.GetSysPthSep() + "config"
	_ = os.Mkdir(configPath, 0764)

	filePath := configPath + commonTool.GetSysPthSep() + "config.xml"
	fileInfo, err := os.ReadFile(filePath)
	if err != nil {
		fileInfo = []byte{}
	}
	return fileInfo
}
