package ui

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/module/server"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"KeyMouseSimulation/share/language"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"math"
	"sort"
	"strconv"
	"time"
)

type ControlT struct {
	wc server.ControlI
	mw *walk.MainWindow

	//记录、回放、暂停、停止按钮
	recordButton      *walk.PushButton
	playbackButton    *walk.PushButton
	pauseButton       *walk.PushButton
	stopButton        *walk.PushButton
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	//热键
	keyList []string
	hKList  [4]string
	hKBox   [4]*walk.ComboBox

	//文件选择
	fileLabel *walk.Label
	fileNames []string
	fileBox   *walk.ComboBox

	//回放次数调整
	playbackTimesLabel *walk.Label
	playbackTimesEdit  *walk.NumberEdit
	currentTimesLabel  *walk.Label
	currentTimesEdit   *walk.NumberEdit

	//速度
	speedLabel *walk.Label
	speedSli   *walk.Slider
	speedEdit  *walk.NumberEdit

	//系统状态
	inFileBox   bool
	statusLabel *walk.Label
	statusEdit  *walk.LineEdit
	errorLabel  *walk.Label
	errorEdit   *walk.TextEdit

	//工具
	settingMenu     *walk.Action
	setHotkeyAction *walk.Action

	languageMenu *walk.Action
	helpMenu     *walk.Action
	aboutAction  *walk.Action
}

var c *ControlT

func createControl() *ControlT {
	c = &ControlT{}
	c.wc = server.NewWinControl()

	//key列表获取
	c.hKList, c.keyList = c.wc.GetKeyList()

	//注册事件订阅
	eventCenter.Event.Register(events.ServerError, c.SubShowError)
	eventCenter.Event.Register(events.ServerHotKeyDown, c.SubHotKeyDown)
	eventCenter.Event.Register(events.ServerChange, c.SubServerChange)

	return c
}

// ----------------------- 主窗口 -----------------------

func MainWindows() {
	c = createControl()
	_, err := MainWindow{
		AssignTo: &c.mw,
		Title:    language.MainWindowTitleStr,
		Size:     Size{Width: 320, Height: 240},
		Layout:   Grid{Columns: 8, Alignment: AlignHNearVCenter},
		Children: []Widget{
			//基础按钮
			PushButton{AssignTo: &c.recordButton, ColumnSpan: 4, Text: language.RecordStr + " " + c.hKList[0], OnClicked: func() {
				c.wc.Record()
			}},
			PushButton{AssignTo: &c.playbackButton, ColumnSpan: 4, Text: language.PlaybackStr + " " + c.hKList[1], OnClicked: func() {
				c.wc.Playback()
			}},
			PushButton{AssignTo: &c.pauseButton, ColumnSpan: 4, Text: language.PauseStr + " " + c.hKList[2], OnClicked: func() {
				_ = c.wc.Pause()
			}},
			PushButton{AssignTo: &c.stopButton, ColumnSpan: 4, Text: language.StopStr + " " + c.hKList[3], OnClicked: func() {
				save := c.wc.Pause()

				//如果需要保存文件
				if save {
					fileName, cmd, err := c.setFileName(c.mw)
					if err != nil {
						_ = c.errorEdit.SetText(err.Error())
						return
					}

					save = cmd == walk.DlgCmdOK //用户点击确认需要保持文件
					if cmd == walk.DlgCmdOK {
						for _, v := range c.fileNames {
							if v == fileName {
								fileName += "-" + strconv.Itoa(int(time.Now().Unix()))
							}
						}
						c.wc.SetFileName(fileName)
						_ = c.fileBox.SetText(fileName)

						c.fileNames = append(c.fileNames, fileName)
						sort.Strings(c.fileNames)
						_ = c.fileBox.SetModel(c.fileNames)
					}
				}

				c.wc.Stop(save)
			}},

			//鼠标路径
			Label{AssignTo: &c.ifMouseTrackLabel, ColumnSpan: 4, Text: language.MouseTrackStr},
			CheckBox{AssignTo: &c.ifMouseTrackCheck, ColumnSpan: 4, Checked: true, Alignment: AlignHCenterVCenter, OnCheckedChanged: func() {
				c.wc.SetIfTrackMouseMove(c.ifMouseTrackCheck.Checked())
			}},

			//回放文件
			Label{AssignTo: &c.fileLabel, Text: language.FileLabelStr, ColumnSpan: 2},
			ComboBox{AssignTo: &c.fileBox, ColumnSpan: 6, Model: c.fileNames, Editable: true, OnCurrentIndexChanged: func() {
				c.wc.SetFileName(c.fileBox.Text())
			}},

			//速度设置
			Label{AssignTo: &c.speedLabel, Text: language.SpeedLabelStr, ColumnSpan: 2},
			NumberEdit{AssignTo: &c.speedEdit, ColumnSpan: 6, MinValue: 0.01, MaxValue: 100.00, Suffix: " X", Decimals: 4},
			Slider{AssignTo: &c.speedSli, ColumnSpan: 8, MinValue: 0, MaxValue: 10, Value: 5, ToolTipsHidden: true, OnValueChanged: func() {
				_ = c.speedEdit.SetValue(math.Pow(2, float64(c.speedSli.Value()-5)))
				c.wc.SetSpeed(c.speedEdit.Value())
			}},

			//回放次数与当前次数
			Label{AssignTo: &c.playbackTimesLabel, Text: language.PlayBackTimesLabelStr, ColumnSpan: 2},
			NumberEdit{AssignTo: &c.playbackTimesEdit, ColumnSpan: 3, MinValue: -1, MaxValue: float64(math.MaxInt64), Value: float64(1), StretchFactor: 0, Suffix: " times", Decimals: 0, OnValueChanged: func() {
				c.wc.SetPlaybackTimes(int(c.playbackTimesEdit.Value()))
			}},
			Label{AssignTo: &c.currentTimesLabel, Text: language.CurrentTimesLabelStr, ColumnSpan: 2},
			NumberEdit{AssignTo: &c.currentTimesEdit, StretchFactor: 0, ColumnSpan: 1},

			//当前状态
			Label{AssignTo: &c.statusLabel, Text: language.StatusLabelStr, ColumnSpan: 2},
			LineEdit{AssignTo: &c.statusEdit, ColumnSpan: 6, ReadOnly: true},

			//错误信息
			Label{AssignTo: &c.errorLabel, Text: language.ErrorLabelStr, ColumnSpan: 2},
			TextEdit{AssignTo: &c.errorEdit, ColumnSpan: 6, ReadOnly: true},
		},

		//工具栏
		MenuItems: []MenuItem{
			Menu{AssignActionTo: &c.settingMenu, Text: "&" + language.MenuSettingStr, Items: []MenuItem{
				Action{AssignTo: &c.setHotkeyAction, Text: language.ActionSetHotKeyStr, OnTriggered: func() {
					if cmd, _ := c.setHotKey(c.mw); cmd == walk.DlgCmdOK {
						c.wc.SetHotKey(enum.HOT_KEY_RECORD_START, c.hKList[0])
						c.wc.SetHotKey(enum.HOT_KEY_PLAYBACK_START, c.hKList[1])
						c.wc.SetHotKey(enum.HOT_KEY_PAUSE, c.hKList[2])
						c.wc.SetHotKey(enum.HOT_KEY_STOP, c.hKList[3])

						_ = c.recordButton.SetText(language.RecordStr + " " + c.hKList[0])
						_ = c.playbackButton.SetText(language.PlaybackStr + " " + c.hKList[1])
						_ = c.pauseButton.SetText(language.PauseStr + " " + c.hKList[2])
						_ = c.stopButton.SetText(language.StopStr + " " + c.hKList[3])
					}
				}},
				Menu{AssignActionTo: &c.languageMenu, Text: language.MenuItemLanguageStr, Items: []MenuItem{
					Action{Text: string(language.English), OnTriggered: func() {
						language.UiChange(language.English)
						language.ServerChange(language.English)
						c.changeLanguage()
						//c.showLanguageBoxAction(c.mw)
					}},
					Action{Text: string(language.Chinese), OnTriggered: func() {
						language.UiChange(language.Chinese)
						language.ServerChange(language.Chinese)
						c.changeLanguage()
						//c.showLanguageBoxAction(c.mw)
					}},
				},
				},
			}},
			Menu{AssignActionTo: &c.helpMenu, Text: language.MenuHelpStr, Items: []MenuItem{
				Action{AssignTo: &c.aboutAction, Text: language.ActionAboutStr, OnTriggered: func() {
					c.showAboutBoxAction(c.mw)
				}},
			}},
		},
	}.Run()
	if err != nil {
		logTool.ErrorAJ(err)
		time.Sleep(5 * time.Second)
	}
}

// ----------------------- 弹窗 -----------------------

func (c *ControlT) setFileName(owner walk.Form) (fileName string, cmd int, err error) {
	//判断是否已经在弹窗中
	if c.inFileBox {
		return fileName, cmd, errors.New("Already In Set File. ")
	}
	defer func() { c.inFileBox = false }()
	c.inFileBox = true

	var nameEdit *walk.LineEdit
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	cmd, err = Dialog{AssignTo: &dlg, Title: language.SetFileWindowTitleStr,
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			TextLabel{Text: language.SetFileLabelStr, ColumnSpan: 4},
			LineEdit{AssignTo: &nameEdit, ColumnSpan: 4, OnTextChanged: func() { fileName = nameEdit.Text() }},
			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.OKStr, OnClicked: func() { dlg.Accept() }},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.CancelStr, OnClicked: func() { dlg.Cancel() }},
		},
	}.Run(owner)

	return
}

func (c *ControlT) setHotKey(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var tmpList = c.hKList

	cmd, err := Dialog{AssignTo: &dlg, Title: language.SetHotKeyWindowTitleStr,
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			Label{Text: language.RecordStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[0], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[0], OnCurrentIndexChanged: func() {
				tmpList[0] = c.hKBox[0].Text()
			}},

			Label{Text: language.PlaybackStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[1], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[1], OnCurrentIndexChanged: func() {
				tmpList[1] = c.hKBox[1].Text()
			}},

			Label{Text: language.PauseStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[2], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[2], OnCurrentIndexChanged: func() {
				tmpList[2] = c.hKBox[2].Text()
			}},

			Label{Text: language.StopStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[3], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[3], OnCurrentIndexChanged: func() {
				tmpList[3] = c.hKBox[3].Text()
			}},

			PushButton{ColumnSpan: 4, Text: language.ResetStr, OnClicked: func() {
				tmpList = [4]string{"F7", "F8", "F9", "F10"}
				_ = c.hKBox[0].SetText(tmpList[0])
				_ = c.hKBox[1].SetText(tmpList[1])
				_ = c.hKBox[2].SetText(tmpList[2])
				_ = c.hKBox[3].SetText(tmpList[3])
			}},

			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: language.OKStr, OnClicked: func() {
				M := make(map[string]bool)
				for _, v := range tmpList {
					if M[v] {
						walk.MsgBox(dlg, language.ErrWindowTitleStr, language.SetHotKeyErrMessageStr, walk.MsgBoxIconInformation)
						return
					} else {
						M[v] = true
					}
				}
				c.hKList = tmpList
				dlg.Accept()
			}},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: language.CancelStr, OnClicked: func() { dlg.Cancel() }},
		},
	}.Run(owner)

	return cmd, err
}

// ----------------------- 工具 -----------------------

func (c *ControlT) showAboutBoxAction(owner walk.Form) {
	walk.MsgBox(owner, language.AboutWindowTitleStr, language.AboutMessageStr, walk.MsgBoxIconInformation)
}

func (c *ControlT) showLanguageBoxAction(owner walk.Form) {
	walk.MsgBox(owner, language.SetLanguageWindowTitleStr, language.SetLanguageChangeMessageStr, walk.MsgBoxIconInformation)
}

// ----------------------- 其他 -----------------------

func (c *ControlT) changeLanguage() {
	_ = c.mw.SetTitle(language.MainWindowTitleStr)

	_ = c.recordButton.SetText(language.RecordStr + " " + c.hKList[0])
	_ = c.playbackButton.SetText(language.PlaybackStr + " " + c.hKList[1])
	_ = c.pauseButton.SetText(language.PauseStr + " " + c.hKList[2])
	_ = c.stopButton.SetText(language.StopStr + " " + c.hKList[3])
	_ = c.ifMouseTrackLabel.SetText(language.MouseTrackStr)
	_ = c.fileLabel.SetText(language.FileLabelStr)
	_ = c.playbackTimesLabel.SetText(language.PlayBackTimesLabelStr)
	_ = c.currentTimesLabel.SetText(language.CurrentTimesLabelStr)
	_ = c.speedLabel.SetText(language.SpeedLabelStr)
	_ = c.statusLabel.SetText(language.StatusLabelStr)
	_ = c.errorLabel.SetText(language.ErrorLabelStr)
	_ = c.mw.SetTitle(language.MainWindowTitleStr)

	_ = c.settingMenu.SetText(language.MenuSettingStr)
	_ = c.languageMenu.SetText(language.MenuHelpStr)
	_ = c.setHotkeyAction.SetText(language.ActionSetHotKeyStr)
	_ = c.helpMenu.SetText(language.MenuHelpStr)
	_ = c.aboutAction.SetText(language.ActionAboutStr)

	c.mw.SetVisible(false)
	c.mw.SetVisible(true)

	time.Sleep(100 * time.Millisecond)
}
func (c *ControlT) waitWidgetLoading() {
	//等待初始化
	for {
		if c.statusEdit != nil && c.playbackTimesEdit != nil && c.errorEdit != nil && c.fileBox != nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (c *ControlT) clickButton(button *walk.PushButton) {
	if c.inFileBox {
		return
	}

	c.mw.WndProc(button.Handle(), win.WM_LBUTTONDOWN, 0, 0)
	c.mw.WndProc(button.Handle(), win.WM_LBUTTONUP, 0, 0)
}

// ----------------------- Sub -----------------------

func (c *ControlT) SubShowError(data interface{}) (err error) {
	d := data.(events.ServerErrorData)

	if c.errorEdit == nil {
		return
	}

	return c.errorEdit.SetText(d.ErrInfo)
}

func (c *ControlT) SubHotKeyDown(data interface{}) (err error) {
	d := data.(events.ServerHotKeyDownData)

	switch d.HotKey {
	case enum.HOT_KEY_PLAYBACK_START:
		go c.clickButton(c.playbackButton)
	case enum.HOT_KEY_RECORD_START:
		go c.clickButton(c.recordButton)
	case enum.HOT_KEY_PAUSE:
		go c.clickButton(c.pauseButton)
	case enum.HOT_KEY_STOP:
		go c.clickButton(c.stopButton)
	}
	return
}
func (c *ControlT) SubServerChange(data interface{}) (err error) {
	d := data.(events.ServerChangeData)

	for c.statusEdit == nil || c.currentTimesEdit == nil || c.fileBox == nil {

	}

	if err = c.statusEdit.SetText(string(d.Status)); err != nil {
		return
	}
	if err = c.currentTimesEdit.SetValue(float64(d.CurrentTimes)); err != nil {
		return
	}
	if d.FileNamesData.Change {
		c.fileNames = d.FileNamesData.FileNames
		if err = c.fileBox.SetModel(d.FileNamesData.FileNames); err != nil {
			return
		}
		if c.fileBox.Text() == "" && len(d.FileNamesData.FileNames) != 0 {
			c.wc.SetFileName(d.FileNamesData.FileNames[0])
			if err = c.fileBox.SetText(d.FileNamesData.FileNames[0]); err != nil {
				return
			}
		}
	}

	return

}
