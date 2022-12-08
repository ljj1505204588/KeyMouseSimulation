package language

type language string

const (
	Chinese language = "中文"
	English          = "English"
)

func init() {
	UiChange(Chinese)
	ServerChange(Chinese)
}

/*
	----------------------------- 界面文字初始化 -----------------------------
*/
var MainWindowTitleStr string
var SetFileWindowTitleStr string
var SetHotKeyWindowTitleStr string
var ErrWindowTitleStr string
var AboutWindowTitleStr string
var SetLanguageWindowTitleStr string

var MenuSettingStr string
var ActionSetHotKeyStr string
var MenuItemLanguageStr string
var MenuHelpStr string
var ActionAboutStr string

var RecordStr string
var PlaybackStr string
var PauseStr string
var StopStr string
var ResetStr string
var MouseTrackStr string
var OKStr string
var CancelStr string

var FileLabelStr string
var SpeedLabelStr string
var PlayBackTimesLabelStr string
var CurrentTimesLabelStr string
var StatusLabelStr string
var ErrorLabelStr string
var SetFileLabelStr string

var AboutMessageStr string
var SetHotKeyErrMessageStr string
var SetLanguageChangeMessageStr string

func UiChange(l language) {
	if l == English {
		MainWindowTitleStr = "RecordAndPlayback"
		SetFileWindowTitleStr = "SetFileName"
		SetHotKeyWindowTitleStr = "SetHotKey"
		ErrWindowTitleStr = "Err"
		AboutWindowTitleStr = "About"
		SetLanguageWindowTitleStr = "Language Change"

		MenuSettingStr = "Setting"
		ActionSetHotKeyStr = "SetHotKey"
		MenuItemLanguageStr = "Language"
		MenuHelpStr = "Help"
		ActionAboutStr = "About"

		RecordStr = "Record"
		PlaybackStr = "Playback"
		PauseStr = "Pause"
		StopStr = "Stop"
		ResetStr = "Reset"
		MouseTrackStr = "Mouse Track"
		OKStr = "OK"
		CancelStr = "Cancel"

		FileLabelStr = "File"
		SpeedLabelStr = "Speed"
		PlayBackTimesLabelStr = "PlaybackTimes"
		CurrentTimesLabelStr = "CurrentTimes"
		StatusLabelStr = "Status"
		ErrorLabelStr = "Error"
		SetFileLabelStr = "File:"

		AboutMessageStr = "MouseKey Record And Playback Tool \n Version:0.1v"
		SetHotKeyErrMessageStr = "Hotkey Settings are repeated"
		SetLanguageChangeMessageStr = "success!"
	} else if l == Chinese {
		MainWindowTitleStr = "鼠标键盘录制回放工具"
		SetFileWindowTitleStr = "文件名称设置"
		SetHotKeyWindowTitleStr = "热键设置"
		ErrWindowTitleStr = "错误"
		AboutWindowTitleStr = "关于 "
		SetLanguageWindowTitleStr = "语言设置"

		MenuSettingStr = "设置"
		ActionSetHotKeyStr = "热键设置"
		MenuItemLanguageStr = "语言设置"
		MenuHelpStr = "帮助"
		ActionAboutStr = "关于"

		RecordStr = "记录"
		PlaybackStr = "回放"
		PauseStr = "暂停"
		StopStr = "停止"
		ResetStr = "重置"
		MouseTrackStr = "鼠标路径"
		OKStr = "确认"
		CancelStr = "取消"

		FileLabelStr = "文件"
		SpeedLabelStr = "速度"
		PlayBackTimesLabelStr = "回放次数"
		CurrentTimesLabelStr = "剩余次数"
		StatusLabelStr = "状态"
		ErrorLabelStr = "错误"
		SetFileLabelStr = "文件名称："

		AboutMessageStr = "鼠标键盘录制回放.\n版本：0.1v"
		SetHotKeyErrMessageStr = "热键设置重复"
		SetLanguageChangeMessageStr = "成功!"
	}
}

/*
	----------------------------- 服务初始化 -----------------------------
*/
var ControlTypeFreeStr string
var ControlTypeRecordingStr string
var ControlTypeRecordPauseStr string
var ControlTypePlaybackStr string
var ControlTypePlaybackPauseStr string

var HotKeyStopStr string
var HotKeyPauseStr string
var HotKeyRecordStartStr string
var HotKeyPlaybackStartStr string

var ErrorPauseFail string
var ErrorFreeToRecordPause string
var ErrorFreeToPlaybackPause string
var ErrorPlaybackToRecordOrRecordPause string
var ErrorPlaybackPauseToRecordOrRecordPause string
var ErrorRecordToPlaybackOrPlaybackPause string
var ErrorRecordPauseToPlaybackOrPlaybackPause string

var ErrorSetHotkeyInFreeStatusStr string
var ErrorSetHotkeyWithoutHotkeyStr string

var ErrorSaveFileNameNilStr string

func ServerChange(l language) {
	if l == English {
		ControlTypeFreeStr = "FREE"
		ControlTypeRecordingStr = "RECORDING"
		ControlTypeRecordPauseStr = "RECORD_PAUSE"
		ControlTypePlaybackStr = "PLAYBACK"
		ControlTypePlaybackPauseStr = "PLAYBACK_PAUSE"

		HotKeyStopStr = "HotKeyStop"
		HotKeyPauseStr = "HotKeyPause"
		HotKeyRecordStartStr = "HotKeyRecordStart"
		HotKeyPlaybackStartStr = "HotKeyPlaybackStart"

		ErrorPauseFail = "Pause Fail. "
		ErrorFreeToRecordPause = "Not in record. "
		ErrorFreeToPlaybackPause = "Not in playback. "
		ErrorPlaybackToRecordOrRecordPause = "Is in playback. please Stop first. "
		ErrorPlaybackPauseToRecordOrRecordPause = "Is in playback. please Stop first. "
		ErrorRecordToPlaybackOrPlaybackPause = "Is in record. please Stop first. "
		ErrorRecordPauseToPlaybackOrPlaybackPause = "Is in record. please Stop first. "
		ErrorSetHotkeyInFreeStatusStr = "Only set hotkey in free status . "
		ErrorSetHotkeyWithoutHotkeyStr = "Hotkey is nil. "
		ErrorSaveFileNameNilStr = "File name is nil. "
	} else if l == Chinese {
		ControlTypeFreeStr = "空闲"
		ControlTypeRecordingStr = "录制"
		ControlTypeRecordPauseStr = "录制暂停"
		ControlTypePlaybackStr = "回放"
		ControlTypePlaybackPauseStr = "回放暂停"

		HotKeyStopStr = "热键停止"
		HotKeyPauseStr = "热键暂停"
		HotKeyRecordStartStr = "热键记录"
		HotKeyPlaybackStartStr = "热键回放"

		ErrorPauseFail = "暂停失败。 "
		ErrorFreeToRecordPause = "当前不在录制中。"
		ErrorFreeToPlaybackPause = "当前不在回放中。"
		ErrorPlaybackToRecordOrRecordPause = "处于回放状态，请先停止。"
		ErrorPlaybackPauseToRecordOrRecordPause = "处于回放状态，请先停止。"
		ErrorRecordToPlaybackOrPlaybackPause = "处于录制状态，请先停止。"
		ErrorRecordPauseToPlaybackOrPlaybackPause = "处于录制状态，请先停止。"
		ErrorSetHotkeyInFreeStatusStr = "只有空闲状态能够设置热键。"
		ErrorSetHotkeyWithoutHotkeyStr = "非法的空值热键。"
		ErrorSaveFileNameNilStr = "文件名称为空。"
	}
}
