package ui

import (
	. "KeyMouseSimulation/module/language"
	"KeyMouseSimulation/module/server"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"math"
	"os"
	"strconv"
	"time"
)

type controlT struct {
	wc server.ControlI
	mw *walk.MainWindow

	//监听
	monitorChan chan server.MessageT

	//记录、回放、暂停、停止按钮
	recordButton      *walk.PushButton
	playbackButton    *walk.PushButton
	pauseButton       *walk.PushButton
	stopButton        *walk.PushButton
	ifMouseTrackLabel *walk.Label
	ifMouseTrackCheck *walk.CheckBox

	//热键
	keyList        []string
	hKList         [4]string
	hKBox          [4]*walk.ComboBox
	recordHKEdit   *walk.LineEdit
	playbackHKEdit *walk.LineEdit
	pauseHKEdit    *walk.LineEdit
	stopHKEdit     *walk.LineEdit

	//文件选择
	fileLabel *walk.Label
	basePath  string
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
	status      int
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

var c *controlT

func createControl() *controlT {
	c = &controlT{}
	c.wc = server.GetWinControl()

	//文件信息获取
	var err error
	c.fileNames = c.wc.ScanFile()
	if c.basePath, err = os.Getwd(); err != nil {
		panic(err.Error())
	}

	//key列表获取
	c.keyList = c.wc.GetKeyList()
	c.hKList = [4]string{"F7", "F8", "F9", "F10"}

	go c.Monitor()
	return c
}
func (c *controlT) Monitor() {
	defer func() {
		if info := recover(); info != nil {
			go c.Monitor()
		}
	}()

	if c.monitorChan == nil {
		c.monitorChan = c.wc.GetMessageChan()
	}

	//等待初始化
	for {
		if c.statusEdit != nil && c.playbackTimesEdit != nil && c.errorEdit != nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	//监听，修改
	for msg := range c.monitorChan {
		switch msg.Event {
		case server.CONTROL_EVENT_PLAYBACK_TIMES_CHANGE:
			if value, ok := msg.Value.(int); ok {
				_ = c.currentTimesEdit.SetValue(float64(value))
			}
		case server.CONTROL_EVENT_STATUS_CHANGE:
			if value, ok := msg.Value.(int); ok {
				c.status = value
			}
			if value2, ok := msg.Value2.(string); ok {
				_ = c.statusEdit.SetText(value2)
			}
		case server.CONTROL_EVENT_ERROR:
			if value, ok := msg.Value.(string); ok {
				_ = c.errorEdit.SetText(value)
			}
		case server.CONTROL_EVENT_SAVE_FILE_ERROR:
			if value, ok := msg.Value.(string); ok {
				_ = c.errorEdit.SetText(value)
			}
		case server.CONTROL_EVENT_HOTKEY_DOWN:
			if value, ok := msg.Value.(server.HotKey); ok {
				switch value {
				case server.HOT_KEY_PLAYBACK_START:
					c.mw.WndProc(c.playbackButton.Handle(), win.WM_LBUTTONDOWN, 0, 0)
					c.mw.WndProc(c.playbackButton.Handle(), win.WM_LBUTTONUP, 0, 0)
				case server.HOT_KEY_RECORD_START:
					c.mw.WndProc(c.recordButton.Handle(), win.WM_LBUTTONDOWN, 0, 0)
					c.mw.WndProc(c.recordButton.Handle(), win.WM_LBUTTONUP, 0, 0)
				case server.HOT_KEY_PUASE:
					c.mw.WndProc(c.pauseButton.Handle(), win.WM_LBUTTONDOWN, 0, 0)
					c.mw.WndProc(c.pauseButton.Handle(), win.WM_LBUTTONUP, 0, 0)
				case server.HOT_KEY_STOP:
					c.mw.WndProc(c.stopButton.Handle(), win.WM_LBUTTONDOWN, 0, 0)
					c.mw.WndProc(c.stopButton.Handle(), win.WM_LBUTTONUP, 0, 0)
				}
			}
		}
	}

	go c.Monitor()
}

// ----------------------- 主窗口 -----------------------

func MainWindows() {
	c = createControl()
	_, err := MainWindow{
		AssignTo: &c.mw,
		Title:    MainWindowTitleStr,
		Size:     Size{Width: 320, Height: 240},
		Layout:   Grid{Columns: 8, Alignment: AlignHNearVCenter},
		Children: []Widget{
			//基础按钮
			PushButton{AssignTo: &c.recordButton, ColumnSpan: 4, Text: RecordStr + " " + c.hKList[0], OnClicked: func() {
				err := c.wc.StartRecord()
				if err != nil {
					_ = c.errorEdit.SetText(err.Error())
				}
			}},
			PushButton{AssignTo: &c.playbackButton, ColumnSpan: 4, Text: PlaybackStr + " " + c.hKList[1], OnClicked: func() {
				err := c.wc.StartPlayback()
				if err != nil {
					_ = c.errorEdit.SetText(err.Error())
				}
			}},
			PushButton{AssignTo: &c.pauseButton, ColumnSpan: 4, Text: PauseStr + " " + c.hKList[2], OnClicked: func() {
				err := c.wc.Pause()
				if err != nil {
					_ = c.errorEdit.SetText(err.Error())
				}
			}},
			PushButton{AssignTo: &c.stopButton, ColumnSpan: 4, Text: StopStr + " " + c.hKList[3], OnClicked: func() {
				//记录中，弹窗
				if c.status == UI_TYPE_RECORDING || c.status == UI_TYPE_RECORD_PAUSE {
					if err := c.wc.Pause(); err != nil {
						_ = c.errorEdit.SetText(err.Error())
					}

					if fileName, cmd, err := c.setFileName(c.mw); err != nil {
						_ = c.errorEdit.SetText(err.Error())
						return
					} else if cmd == walk.DlgCmdOK {
						for _, v := range c.fileNames {
							if v == fileName {
								fileName += "-" + strconv.Itoa(int(time.Now().Unix()))
							}
						}
						c.wc.SetFileName(fileName)
						_ = c.fileBox.SetText(fileName)
						c.fileNames = append(c.fileNames, fileName)
						_ = c.fileBox.SetModel(c.fileNames)
					}
				}

				if err := c.wc.Stop(); err != nil {
					_ = c.errorEdit.SetText(err.Error())
					return
				}
			}},

			//鼠标路径
			Label{AssignTo: &c.ifMouseTrackLabel,ColumnSpan: 4,Text:MouseTrackStr },
			CheckBox{AssignTo: &c.ifMouseTrackCheck,ColumnSpan: 4,Checked: true,Alignment:AlignHCenterVCenter,OnCheckedChanged: func() {
				c.wc.SetIfTrackMouseMove(c.ifMouseTrackCheck.Checked())
			}},

			//回放文件
			Label{AssignTo: &c.fileLabel, Text: FileLabelStr, ColumnSpan: 2},
			ComboBox{AssignTo: &c.fileBox, ColumnSpan: 6, Model: c.fileNames, Editable: true, OnCurrentIndexChanged: func() {
				c.wc.SetFileName(c.fileBox.Text())
			}},

			//速度设置
			Label{AssignTo: &c.speedLabel, Text: SpeedLabelStr, ColumnSpan: 2},
			NumberEdit{AssignTo: &c.speedEdit, ColumnSpan: 6, MinValue: 0.01, MaxValue: 100.00, Suffix: " X", Decimals: 4},
			Slider{AssignTo: &c.speedSli, ColumnSpan: 8, MinValue: 0, MaxValue: 10, Value: 5, ToolTipsHidden: true, OnValueChanged: func() {
				_ = c.speedEdit.SetValue(math.Pow(2, float64(c.speedSli.Value()-5)))
				c.wc.SetSpeed(c.speedEdit.Value())
			}},

			//回放次数与当前次数
			Label{AssignTo: &c.playbackTimesLabel, Text: PlayBackTimesLabelStr, ColumnSpan: 2},
			NumberEdit{AssignTo: &c.playbackTimesEdit, ColumnSpan: 3, MinValue: -1, MaxValue: float64(math.MaxInt64), StretchFactor: 0, Suffix: " times", Decimals: 0, OnValueChanged: func() {
				c.wc.SetPlaybackTimes(int(c.playbackTimesEdit.Value()))
			}},
			Label{AssignTo: &c.currentTimesLabel, Text: CurrentTimesLabelStr, ColumnSpan: 2},
			NumberEdit{AssignTo: &c.currentTimesEdit, StretchFactor: 0, ColumnSpan: 1},

			//当前状态
			Label{AssignTo: &c.statusLabel, Text: StatusLabelStr, ColumnSpan: 2},
			LineEdit{AssignTo: &c.statusEdit, ColumnSpan: 6, ReadOnly: true},

			//错误信息
			Label{AssignTo: &c.errorLabel, Text: ErrorLabelStr, ColumnSpan: 2},
			TextEdit{AssignTo: &c.errorEdit, ColumnSpan: 6, ReadOnly: true},
		},

		//工具栏
		MenuItems: []MenuItem{
			Menu{AssignActionTo: &c.settingMenu, Text: "&" + MenuSettingStr, Items: []MenuItem{
				Action{AssignTo: &c.setHotkeyAction, Text: ActionSetHotKeyStr, OnTriggered: func() {
					if cmd, _ := c.setHotKey(c.mw); cmd == walk.DlgCmdOK {
						c.wc.SetHotKey(server.HOT_KEY_RECORD_START, c.hKList[0])
						c.wc.SetHotKey(server.HOT_KEY_PLAYBACK_START, c.hKList[1])
						c.wc.SetHotKey(server.HOT_KEY_PUASE, c.hKList[2])
						c.wc.SetHotKey(server.HOT_KEY_STOP, c.hKList[3])

						_ = c.recordButton.SetText(RecordStr + " " + c.hKList[0])
						_ = c.playbackButton.SetText(PlaybackStr + " " + c.hKList[1])
						_ = c.pauseButton.SetText(PauseStr + " " + c.hKList[2])
						_ = c.stopButton.SetText(StopStr + " " + c.hKList[3])
					}
				}},
				Menu{AssignActionTo: &c.languageMenu, Text: MenuItemLanguageStr, Items: []MenuItem{
					Action{Text: string(English), OnTriggered: func() {
						UiChange(English)
						ServerChange(English)
						c.ChangeLanguage()
						//c.showLanguageBoxAction(c.mw)
					}},
					Action{Text: string(Chinese), OnTriggered: func() {
						UiChange(Chinese)
						ServerChange(Chinese)
						c.ChangeLanguage()
						//c.showLanguageBoxAction(c.mw)
					}},
				},
				},
			}},
			Menu{AssignActionTo: &c.helpMenu, Text: MenuHelpStr, Items: []MenuItem{
				Action{AssignTo: &c.aboutAction, Text: ActionAboutStr, OnTriggered: func() {
					c.showAboutBoxAction(c.mw)
				}},
			}},
		},
	}.Run()
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(5*time.Second)
	}
}

// ----------------------- 弹窗 -----------------------

func (c *controlT) setFileName(owner walk.Form) (string, int, error) {
	var nameEdit *walk.LineEdit
	filename := ""
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	cmd, err := Dialog{AssignTo: &dlg, Title: SetFileWindowTitleStr,
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			TextLabel{Text: SetFileLabelStr, ColumnSpan: 4},
			LineEdit{AssignTo: &nameEdit, ColumnSpan: 4, OnTextChanged: func() { filename = nameEdit.Text() }},
			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: OKStr, OnClicked: func() { dlg.Accept() }},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: CancelStr, OnClicked: func() { dlg.Cancel() }},
		},
	}.Run(owner)

	return filename, cmd, err
}

func (c *controlT) setHotKey(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var tmpList = c.hKList

	cmd, err := Dialog{AssignTo: &dlg, Title: SetHotKeyWindowTitleStr,
		DefaultButton: &acceptPB, CancelButton: &cancelPB,
		Size: Size{Width: 350, Height: 200}, Layout: Grid{Columns: 4},
		Children: []Widget{
			Label{Text: RecordStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[0], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[0], OnCurrentIndexChanged: func() {
				tmpList[0] = c.hKBox[0].Text()
			}},

			Label{Text: PlaybackStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[1], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[1], OnCurrentIndexChanged: func() {
				tmpList[1] = c.hKBox[1].Text()
			}},

			Label{Text: PauseStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[2], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[2], OnCurrentIndexChanged: func() {
				tmpList[2] = c.hKBox[2].Text()
			}},

			Label{Text: StopStr, ColumnSpan: 1},
			ComboBox{AssignTo: &c.hKBox[3], ColumnSpan: 1, Model: c.keyList, Editable: true, Value: c.hKList[3], OnCurrentIndexChanged: func() {
				tmpList[3] = c.hKBox[3].Text()
			}},

			PushButton{ColumnSpan: 4, Text: ResetStr, OnClicked: func() {
				tmpList = [4]string{"F7", "F8", "F9", "F10"}
				_ = c.hKBox[0].SetText(tmpList[0])
				_ = c.hKBox[1].SetText(tmpList[1])
				_ = c.hKBox[2].SetText(tmpList[2])
				_ = c.hKBox[3].SetText(tmpList[3])
			}},

			PushButton{AssignTo: &acceptPB, ColumnSpan: 2, Text: OKStr, OnClicked: func() {
				M := make(map[string]bool)
				for _, v := range tmpList {
					if M[v] {
						walk.MsgBox(dlg, ErrWindowTitleStr, SetHotKeyErrMessageStr, walk.MsgBoxIconInformation)
						return
					} else {
						M[v] = true
					}
				}
				c.hKList = tmpList
				dlg.Accept()
			}},
			PushButton{AssignTo: &cancelPB, ColumnSpan: 2, Text: CancelStr, OnClicked: func() { dlg.Cancel() }},
		},
	}.Run(owner)

	return cmd, err
}

// ----------------------- 工具 -----------------------

func (c *controlT) showAboutBoxAction(owner walk.Form) {
	walk.MsgBox(owner, AboutWindowTitleStr, AboutMessageStr, walk.MsgBoxIconInformation)
}

func (c *controlT) showLanguageBoxAction(owner walk.Form) {
	walk.MsgBox(owner, SetLanguageWindowTitleStr, SetLanguageChangeMessageStr, walk.MsgBoxIconInformation)
}

// ----------------------- 其他 -----------------------

func (c *controlT) ChangeLanguage() {
	_ = c.mw.SetTitle(MainWindowTitleStr)

	_ = c.recordButton.SetText(RecordStr + " " + c.hKList[0])
	_ = c.playbackButton.SetText(PlaybackStr + " " + c.hKList[1])
	_ = c.pauseButton.SetText(PauseStr + " " + c.hKList[2])
	_ = c.stopButton.SetText(StopStr + " " + c.hKList[3])
	_ = c.ifMouseTrackLabel.SetText(MouseTrackStr)
	_ = c.fileLabel.SetText(FileLabelStr)
	_ = c.playbackTimesLabel.SetText(PlayBackTimesLabelStr)
	_ = c.currentTimesLabel.SetText(CurrentTimesLabelStr)
	_ = c.speedLabel.SetText(SpeedLabelStr)
	_ = c.statusLabel.SetText(StatusLabelStr)
	_ = c.errorLabel.SetText(ErrorLabelStr)
	_ = c.mw.SetTitle(MainWindowTitleStr)

	_ = c.settingMenu.SetText(MenuSettingStr)
	_ = c.languageMenu.SetText(MenuHelpStr)
	_ = c.setHotkeyAction.SetText(ActionSetHotKeyStr)
	_ = c.helpMenu.SetText(MenuHelpStr)
	_ = c.aboutAction.SetText(ActionAboutStr)

	c.mw.SetVisible(false)
	c.mw.SetVisible(true)

	time.Sleep(100 * time.Millisecond)
}
