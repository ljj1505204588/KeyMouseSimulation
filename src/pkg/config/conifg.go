package conf

import (
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/topic"
	"encoding/json"
	"os"
	"sync/atomic"
)

var (
	path   = "./config/config.json"
	conf   configDataT
	saving int64
)

type configDataT struct {
	RecordHotKeyConf    string `json:"recordHotKeyConf"`
	PlayBackHotKeyConf  string `json:"playBackHotKeyConf"`
	PauseHotKeyConf     string `json:"pauseHotKeyConf"`
	StopHotKeyConf      string `json:"stopHotKeyConf"`
	SpeedUpHotKeyConf   string `json:"speedUpHotKeyConf"`
	SpeedDownHotKeyConf string `json:"speedDownHotKeyConf"`

	RecordMouseTrackConf *bool    `json:"recordMouseTrackConf"`
	PlaybackSpeedConf    *float64 `json:"playbackSpeedConf"`
	PlaybackTimesConf    *int64   `json:"playbackTimesConf"`
	LanguageConf         string   `json:"languageConf"`
}

// 初始化配置
func init() {
	// 确保配置目录存在
	_ = os.MkdirAll("./config", 0755)
	// 加载配置
	loadConfig()
	// 监听配置变更事件
	eventCenter.Event.Register(topic.ConfigChange, handleConfigChange)
}

// 加载配置
func loadConfig() {

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		// 如果文件不存在，使用默认配置
		return
	}

	// 解析JSON
	if err := json.Unmarshal(data, &conf); err != nil {
		return
	}

	// 应用配置到各个配置项
	applyConfig(conf)
}

// 保存配置
func saveConfig() {

	// 转换为JSON
	data, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return
	}

	// 覆盖文件
	var file *os.File
	if file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return
	}
	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return
	}

	return
}

// 应用配置
func applyConfig(conf configDataT) {
	if conf.RecordHotKeyConf != "" {
		RecordHotKeyConf.SetValue(conf.RecordHotKeyConf)
	}
	if conf.PlayBackHotKeyConf != "" {
		PlayBackHotKeyConf.SetValue(conf.PlayBackHotKeyConf)
	}
	if conf.PauseHotKeyConf != "" {
		PauseHotKeyConf.SetValue(conf.PauseHotKeyConf)
	}
	if conf.StopHotKeyConf != "" {
		StopHotKeyConf.SetValue(conf.StopHotKeyConf)
	}
	if conf.SpeedUpHotKeyConf != "" {
		SpeedUpHotKeyConf.SetValue(conf.SpeedUpHotKeyConf)
	}
	if conf.SpeedDownHotKeyConf != "" {
		SpeedDownHotKeyConf.SetValue(conf.SpeedDownHotKeyConf)
	}
	if conf.RecordMouseTrackConf != nil {
		RecordMouseTrackConf.SetValue(*conf.RecordMouseTrackConf)
	}
	if conf.PlaybackSpeedConf != nil {
		PlaybackSpeedConf.SetValue(*conf.PlaybackSpeedConf)
	}
	if conf.PlaybackTimesConf != nil {
		PlaybackTimesConf.SetValue(*conf.PlaybackTimesConf)
	}
	if conf.LanguageConf != "" {
		LanguageConf.SetValue(conf.LanguageConf)
	}
}

// 处理配置变更事件
func handleConfigChange(data interface{}) (err error) {
	//todo 修改下配置修改的逻辑
	if !atomic.CompareAndSwapInt64(&saving, 0, 1) {
		return
	}
	defer atomic.StoreInt64(&saving, 0)

	var (
		RecordMouseTrackValue = RecordMouseTrackConf.GetValue()
		PlaybackSpeedValue    = PlaybackSpeedConf.GetValue()
		PlaybackTimesValue    = PlaybackTimesConf.GetValue()
	)
	conf = configDataT{
		RecordHotKeyConf:     RecordHotKeyConf.GetValue(),
		PlayBackHotKeyConf:   PlayBackHotKeyConf.GetValue(),
		PauseHotKeyConf:      PauseHotKeyConf.GetValue(),
		StopHotKeyConf:       StopHotKeyConf.GetValue(),
		SpeedUpHotKeyConf:    SpeedUpHotKeyConf.GetValue(),
		SpeedDownHotKeyConf:  SpeedDownHotKeyConf.GetValue(),
		RecordMouseTrackConf: &RecordMouseTrackValue,
		PlaybackSpeedConf:    &PlaybackSpeedValue,
		PlaybackTimesConf:    &PlaybackTimesValue,
		LanguageConf:         LanguageConf.GetValue(),
	}

	saveConfig()
	return
}
