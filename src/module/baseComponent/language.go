package component

type (
	DisPlay int64
	Type    string
)

var Center = &centerT{
	languageMap: chinesMap,
}

type centerT struct {
	languageMap   map[DisPlay]string
	changeHandler []func()
}

func (c *centerT) Get(no DisPlay) string {
	return c.languageMap[no]
}
func (c *centerT) SetLanguage(typ Type) {
	switch typ {
	case Chinese:
		c.languageMap = chinesMap
	case English:
		c.languageMap = englishMap
	}

	for _, h := range c.changeHandler {
		h()
	}
}

func (c *centerT) Refresh() {
	for _, h := range c.changeHandler {
		h()
	}
}

func (c *centerT) RegisterChange(h ...func()) {
	c.changeHandler = append(c.changeHandler, h...)
}

/*
	----------------------------- 界面文字初始化 -----------------------------
*/

const (
	Chinese Type = "中文"
	English Type = "english"
)

const (
	MainWindowTitleStr        DisPlay = iota // 鼠标键盘录制回放工具
	SetFileWindowTitleStr                    // 文件名称设置
	SetHotKeyWindowTitleStr                  // 热键设置
	ErrWindowTitleStr                        // 错误
	AboutWindowTitleStr                      // 关于
	SetLanguageWindowTitleStr                // 语言设置
	MenuSettingStr                           // 设置
	ActionSetHotKeyStr                       // 热键设置
	MenuItemLanguageStr                      // 语言设置
	MenuHelpStr                              // 帮助
	ActionAboutStr                           // 关于
	RecordStr                                // 记录
	PlaybackStr                              // 回放
	PauseStr                                 // 暂停
	StopStr                                  // 停止
	ResetStr                                 // 重置
	MouseTrackStr                            // 鼠标路径
	OKStr                                    // 确认
	CancelStr                                // 取消
	FileLabelStr                             // 文件
	SpeedLabelStr                            // 速度
	PlayBackTimesLabelStr                    // 回放次数
	CurrentTimesLabelStr                     // 剩余次数
	StatusLabelStr                           // 状态
	ErrorLabelStr                            // 错误

	SetFileLabelStr             // 文件名称：
	AboutMessageStr             // 鼠标键盘录制回放.\n版本：0.1v
	SetHotKeyErrMessageStr      // 热键设置重复
	SetLanguageChangeMessageStr // 成功!

	ControlTypeFreeStr          // 空闲
	ControlTypeRecordingStr     // 录制
	ControlTypeRecordPauseStr   // 录制暂停
	ControlTypePlaybackStr      // 回放
	ControlTypePlaybackPauseStr // 回放暂停

	HotKeyStopStr          // 热键停止
	HotKeyPauseStr         // 热键暂停
	HotKeyRecordStartStr   // 热键记录
	HotKeyPlaybackStartStr // 热键回放

	ErrorStatusChangeError         // 状态变动错误。
	ErrorSetHotkeyInFreeStatusStr  // 只有空闲状态能够设置热键。
	ErrorSetHotkeyWithoutHotkeyStr // 非法的空值热键。
	ErrorSaveFileNameNilStr        // 文件名称为空。
	ErrorSaveFileNotExist          // 文件不存在。
)

var chinesMap = map[DisPlay]string{
	MainWindowTitleStr:          "鼠标键盘录制回放工具",
	SetFileWindowTitleStr:       "文件名称设置",
	SetHotKeyWindowTitleStr:     "热键设置",
	ErrWindowTitleStr:           "错误",
	AboutWindowTitleStr:         "关于 ",
	SetLanguageWindowTitleStr:   "语言设置",
	MenuSettingStr:              "设置",
	ActionSetHotKeyStr:          "热键设置",
	MenuItemLanguageStr:         "语言设置",
	MenuHelpStr:                 "帮助",
	ActionAboutStr:              "关于",
	RecordStr:                   "记录",
	PlaybackStr:                 "回放",
	PauseStr:                    "暂停",
	StopStr:                     "停止",
	ResetStr:                    "重置",
	MouseTrackStr:               "鼠标路径",
	OKStr:                       "确认",
	CancelStr:                   "取消",
	FileLabelStr:                "文件",
	SpeedLabelStr:               "速度",
	PlayBackTimesLabelStr:       "回放次数",
	CurrentTimesLabelStr:        "剩余次数",
	StatusLabelStr:              "状态",
	ErrorLabelStr:               "错误",
	SetFileLabelStr:             "文件名称：",
	AboutMessageStr:             "鼠标键盘录制回放.\n版本：0.1v",
	SetHotKeyErrMessageStr:      "热键设置重复",
	SetLanguageChangeMessageStr: "成功!",

	ControlTypeFreeStr:          "空闲",
	ControlTypeRecordingStr:     "录制",
	ControlTypeRecordPauseStr:   "录制暂停",
	ControlTypePlaybackStr:      "回放",
	ControlTypePlaybackPauseStr: "回放暂停",
	HotKeyStopStr:               "热键停止",
	HotKeyPauseStr:              "热键暂停",
	HotKeyRecordStartStr:        "热键记录",
	HotKeyPlaybackStartStr:      "热键回放",
	ErrorStatusChangeError:      "状态错误调用方法。",

	ErrorSetHotkeyInFreeStatusStr:  "只有空闲状态能够设置热键。",
	ErrorSetHotkeyWithoutHotkeyStr: "非法的空值热键。",
	ErrorSaveFileNameNilStr:        "文件名称为空。",
	ErrorSaveFileNotExist:          "文件不存在。",
}
var englishMap = map[DisPlay]string{
	MainWindowTitleStr:          "RecordAndPlayback",
	SetFileWindowTitleStr:       "SetFileName",
	SetHotKeyWindowTitleStr:     "SetHotKey",
	ErrWindowTitleStr:           "Err",
	AboutWindowTitleStr:         "About",
	SetLanguageWindowTitleStr:   "Language Change",
	MenuSettingStr:              "Setting",
	ActionSetHotKeyStr:          "SetHotKey",
	MenuItemLanguageStr:         "Language",
	MenuHelpStr:                 "Help",
	ActionAboutStr:              "About",
	RecordStr:                   "Record",
	PlaybackStr:                 "Playback",
	PauseStr:                    "Pause",
	StopStr:                     "Stop",
	ResetStr:                    "Reset",
	MouseTrackStr:               "Mouse Track",
	OKStr:                       "OK",
	CancelStr:                   "Cancel",
	FileLabelStr:                "File",
	SpeedLabelStr:               "Speed",
	PlayBackTimesLabelStr:       "PlaybackTimes",
	CurrentTimesLabelStr:        "CurrentTimes",
	StatusLabelStr:              "Status",
	ErrorLabelStr:               "Error",
	SetFileLabelStr:             "File:",
	AboutMessageStr:             "MouseKey Record And Playback Tool \n Version:0.1v",
	SetHotKeyErrMessageStr:      "Hotkey Settings are repeated",
	SetLanguageChangeMessageStr: "success!",

	ControlTypeFreeStr:          "FREE",
	ControlTypeRecordingStr:     "RECORDING",
	ControlTypeRecordPauseStr:   "RECORD_PAUSE",
	ControlTypePlaybackStr:      "PLAYBACK",
	ControlTypePlaybackPauseStr: "PLAYBACK_PAUSE",
	HotKeyStopStr:               "HotKeyStop",
	HotKeyPauseStr:              "HotKeyPause",
	HotKeyRecordStartStr:        "HotKeyRecordStart",
	HotKeyPlaybackStartStr:      "HotKeyPlaybackStart",
	ErrorStatusChangeError:      "ChangeStatusError. ",

	ErrorSetHotkeyInFreeStatusStr:  "Only set hotkey in free status . ",
	ErrorSetHotkeyWithoutHotkeyStr: "Hotkey is nil. ",
	ErrorSaveFileNameNilStr:        "File name is nil. ",
	ErrorSaveFileNotExist:          "File not exist. ",
}
