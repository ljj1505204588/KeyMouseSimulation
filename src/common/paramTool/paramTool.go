package paramTool

import (
	"KeyMouseSimulation/common/commonTool"
	"encoding/xml"
	"fmt"
	"os"
)

const PARAM_TOOL_LOG_HEAD = "配置读取工具："

func init() {
	center.create()

	//获取配置
	text := getConfig()

	//读取配置
	initConfig(text)

	//复写配置回文件
	if err := WriteBackFile(); err != nil {
		panic(PARAM_TOOL_LOG_HEAD + "复写回文件配置编码 Xml 失败.")
	}

	fmt.Println(PARAM_TOOL_LOG_HEAD + "参数组件，加载成功！")
}

//WriteBackFile 回写回文件
func WriteBackFile() error {
	//获取回写内容
	writeBack, err := xml.MarshalIndent(center.getWriterBack(), "", "    ")
	if err != nil {
		return err
	}

	//创建||打开文件
	file, openErr := os.OpenFile(getConfigPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0762)
	if openErr != nil {
		return openErr
	}

	//回写回文件
	_, err = file.Write(writeBack)
	_ = file.Close()
	return err
}

// 获取配置文件
func getConfig() (text []byte) {
	text, _ = os.ReadFile(getConfigPath())
	return
}

// 初始化责任链 读取相关配置
func initConfig(text []byte) {
	center.initParam(text)

}

// 获取配置文件地址
func getConfigPath() string {
	//取软件执行地址
	wd, _ := os.Getwd()

	//生成Config文件夹
	configPath := wd + commonTool.GetSysPthSep() + "config"
	_ = os.Mkdir(configPath, 0764)

	//返回配置地址
	return configPath + commonTool.GetSysPthSep() + "config.xml"
}
