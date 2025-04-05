package language

import (
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
)

var langType = enum.LanguageType(conf.LanguageConf.GetValue()) // todo 改成conf

func init() {
	eventCenter.Event.Register(topic.LanguageChange, func(data interface{}) (err error) {
		var dataValue = data.(*topic.LanguageChangeData)
		langType = dataValue.Typ
		conf.LanguageConf.SetValue(string(dataValue.Typ))
		return
	}, eventCenter.SetOrderLv(100))

}

// -------------------- 语言键 --------------------

// LanguageKeyI 语言键接口
type LanguageKeyI interface {
	ToString() string
}

type languageKeyT struct {
	chinese string // 中文
	english string // 英文
}

// ToString 获取字符串
func (l *languageKeyT) ToString() string {
	if langType == enum.Chinese {
		return l.chinese
	}
	return l.english
}

var (
	MainWindowTitleStr             LanguageKeyI = &languageKeyT{chinese: "鼠标键盘录制回放工具", english: "RecordAndPlayback"}
	SetFileWindowTitleStr          LanguageKeyI = &languageKeyT{chinese: "文件名称设置", english: "SetFileName"}
	SetHotKeyWindowTitleStr        LanguageKeyI = &languageKeyT{chinese: "热键设置", english: "SetHotKey"}
	ErrWindowTitleStr              LanguageKeyI = &languageKeyT{chinese: "错误", english: "Err"}
	AboutWindowTitleStr            LanguageKeyI = &languageKeyT{chinese: "关于 ", english: "About"}
	SetLanguageWindowTitleStr      LanguageKeyI = &languageKeyT{chinese: "语言设置", english: "Language Change"}
	MenuSettingStr                 LanguageKeyI = &languageKeyT{chinese: "设置", english: "Setting"}
	ActionSetHotKeyStr             LanguageKeyI = &languageKeyT{chinese: "热键设置", english: "SetHotKey"}
	MenuItemLanguageStr            LanguageKeyI = &languageKeyT{chinese: "语言设置", english: "Language"}
	MenuHelpStr                    LanguageKeyI = &languageKeyT{chinese: "帮助", english: "Help"}
	ActionAboutStr                 LanguageKeyI = &languageKeyT{chinese: "关于", english: "About"}
	RecordStr                      LanguageKeyI = &languageKeyT{chinese: "记录", english: "Record"}
	PlayBackStr                    LanguageKeyI = &languageKeyT{chinese: "回放", english: "Playback"}
	PauseStr                       LanguageKeyI = &languageKeyT{chinese: "暂停", english: "Pause"}
	StopStr                        LanguageKeyI = &languageKeyT{chinese: "停止", english: "Stop"}
	SpeedUpStr                     LanguageKeyI = &languageKeyT{chinese: "加速", english: "SpeedUp"}
	SpeedDownStr                   LanguageKeyI = &languageKeyT{chinese: "减速", english: "SpeedDown"}
	ResetStr                       LanguageKeyI = &languageKeyT{chinese: "重置", english: "Reset"}
	MouseTrackStr                  LanguageKeyI = &languageKeyT{chinese: "鼠标路径", english: "Mouse Track"}
	RecordLenStr                   LanguageKeyI = &languageKeyT{chinese: "记录长度", english: "RecordLen"}
	OKStr                          LanguageKeyI = &languageKeyT{chinese: "确认", english: "OK"}
	CancelStr                      LanguageKeyI = &languageKeyT{chinese: "取消", english: "Cancel"}
	FileLabelStr                   LanguageKeyI = &languageKeyT{chinese: "文件", english: "File"}
	SpeedLabelStr                  LanguageKeyI = &languageKeyT{chinese: "速度", english: "Speed"}
	PlayBackTimesLabelStr          LanguageKeyI = &languageKeyT{chinese: "回放次数", english: "PlaybackTimes"}
	CurrentTimesLabelStr           LanguageKeyI = &languageKeyT{chinese: "剩余次数", english: "CurrentTimes"}
	StatusLabelStr                 LanguageKeyI = &languageKeyT{chinese: "状态", english: "Status"}
	ErrorLabelStr                  LanguageKeyI = &languageKeyT{chinese: "错误信息", english: "Error"}
	SetFileLabelStr                LanguageKeyI = &languageKeyT{chinese: "文件名称：", english: "File:"}
	AboutMessageStr                LanguageKeyI = &languageKeyT{chinese: "鼠标键盘录制回放.\n版本：1.0v", english: "MouseKey Record And Playback Tool \n Version:1.0v"}
	SetHotKeyErrMessageStr         LanguageKeyI = &languageKeyT{chinese: "热键设置重复", english: "Hotkey Settings are repeated"}
	SetLanguageChangeMessageStr    LanguageKeyI = &languageKeyT{chinese: "成功!", english: "success!"}
	ControlTypeFreeStr             LanguageKeyI = &languageKeyT{chinese: "空闲", english: "FREE"}
	ControlTypeRecordingStr        LanguageKeyI = &languageKeyT{chinese: "录制", english: "RECORDING"}
	ControlTypeRecordPauseStr      LanguageKeyI = &languageKeyT{chinese: "录制暂停", english: "RECORD_PAUSE"}
	ControlTypePlaybackStr         LanguageKeyI = &languageKeyT{chinese: "回放", english: "PLAYBACK"}
	ControlTypePlaybackPauseStr    LanguageKeyI = &languageKeyT{chinese: "回放暂停", english: "PLAYBACK_PAUSE"}
	HotKeyStopStr                  LanguageKeyI = &languageKeyT{chinese: "热键停止", english: "HotKeyStop"}
	HotKeyPauseStr                 LanguageKeyI = &languageKeyT{chinese: "热键暂停", english: "HotKeyPause"}
	HotKeyRecordStartStr           LanguageKeyI = &languageKeyT{chinese: "热键记录", english: "HotKeyRecordStart"}
	HotKeyPlaybackStartStr         LanguageKeyI = &languageKeyT{chinese: "热键回放", english: "HotKeyPlaybackStart"}
	ErrorStatusChangeError         LanguageKeyI = &languageKeyT{chinese: "状态错误调用方法。", english: "ChangeStatusError. "}
	ErrorSetHotkeyInFreeStatusStr  LanguageKeyI = &languageKeyT{chinese: "只有空闲状态能够设置热键。", english: "Only set hotkey in free status . "}
	ErrorSetHotkeyWithoutHotkeyStr LanguageKeyI = &languageKeyT{chinese: "非法的空值热键。", english: "Hotkey is nil. "}
	ErrorSaveFileNameNilStr        LanguageKeyI = &languageKeyT{chinese: "文件名称为空。", english: "File name is nil. "}
	ErrorSaveFileNotExist          LanguageKeyI = &languageKeyT{chinese: "文件不存在。", english: "File not exist. "}
	ErrorTopicUnregistered         LanguageKeyI = &languageKeyT{chinese: "主题未注册。", english: "Topic Unregistered."}
)
