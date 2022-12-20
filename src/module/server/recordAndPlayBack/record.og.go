package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type RecordServerI interface {
	Start()                      //开始
	Pause()                      //暂停
	Stop(name string, save bool) //停止

	SetHotKey(k enum.HotKey, vks keyMouTool.VKCode) error //设置热键
	SetIfTrackMouseMove(sign bool)                        //设置是否记录鼠标移动路径
}

/*
*	RecordServerI 实现接口
 */

func GetRecordServer() *RecordServerT {
	R := RecordServerT{
		hotKeyM:          make(map[keyMouTool.VKCode]enum.HotKey),
		recordMouseTrack: true,
	}

	R.mouseDwMap = map[windowsHook.Message]keyMouTool.MouseInputDW{
		windowsHook.WM_MOUSEMOVE:         keyMouTool.DW_MOUSEEVENTF_MOVE | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTDOWN:     keyMouTool.DW_MOUSEEVENTF_LEFTDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTUP:       keyMouTool.DW_MOUSEEVENTF_LEFTUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSERIGHTDOWN:    keyMouTool.DW_MOUSEEVENTF_RIGHTDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSERIGHTUP:      keyMouTool.DW_MOUSEEVENTF_RIGHTUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEMIDDLEDOWN:   keyMouTool.DW_MOUSEEVENTF_MIDDLEDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEMIDDLEUP:     keyMouTool.DW_MOUSEEVENTF_MIDDLEUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTSICEDOWN: keyMouTool.DW_MOUSEEVENTF_XDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTSICEUP:   keyMouTool.DW_MOUSEEVENTF_XUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEWHEEL:        keyMouTool.DW_MOUSEEVENTF_WHEEL | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEHWHEEL:       keyMouTool.DW_MOUSEEVENTF_HWHEEL | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE, //这个暂时不知道是啥，
	}
	R.keyDwMap = map[windowsHook.Message]keyMouTool.KeyBoardInputDW{
		windowsHook.WM_KEYUP: keyMouTool.DW_KEYEVENTF_KEYUP,
	}

	//初始化监听热键
	R.initHook()

	go R.mainLoop()

	return &R
}

type RecordServerT struct {
	recordStartTime int64   //记录开始时间，纳秒
	notes           []noteT //记录

	hotKeyM          map[keyMouTool.VKCode]enum.HotKey //热键信息
	recordMouseTrack bool                              //是否记录鼠标移动路径使用
	lastMoveEven     *windowsHook.MouseEvent           //最后移动事件，配合是否记录鼠标移动路径使用

	mouseChan    chan windowsHook.MouseEvent    //鼠标监听通道
	keyboardChan chan windowsHook.KeyboardEvent //键盘监听通道
	hotKeyChan   chan windowsHook.KeyboardEvent //热键监听通道

	mouseDwMap map[windowsHook.Message]keyMouTool.MouseInputDW    //鼠标转换Map
	keyDwMap   map[windowsHook.Message]keyMouTool.KeyBoardInputDW //键盘转换Map
}

// Start 开始
func (R *RecordServerT) Start() {
	if R.recordStartTime == 0 {
		R.recordStartTime = time.Now().UnixNano()
		R.handUpHook()
	}
}

// Pause 暂停
func (R *RecordServerT) Pause() {
	R.recordStartTime = 0
	R.handOutHook()
}

// Stop 停止
func (R *RecordServerT) Stop(name string, save bool) {
	R.recordStartTime = 0
	R.handOutHook()

	//记录文件
	notes := R.notes
	R.notes = []noteT{}
	if len(notes) != 0 && save {
		go R.recordNoteToFile(name, notes)
	}

	return
}

// SetHotKey 设置热键
func (R *RecordServerT) SetHotKey(k enum.HotKey, vks keyMouTool.VKCode) error {
	M := make(map[keyMouTool.VKCode]enum.HotKey)

	if R.hotKeyM != nil {
		for k, v := range R.hotKeyM {
			M[k] = v
		}
	}

	M[vks] = k
	R.hotKeyM = M

	return nil
}

// SetIfTrackMouseMove 设置是否追踪鼠标记录
func (R *RecordServerT) SetIfTrackMouseMove(sign bool) {
	R.recordMouseTrack = sign
}

// ----------------------- record 模块主体功能 -----------------------

func (R *RecordServerT) mainLoop() {
	defer func() {
		if info := recover(); info != nil {
			R.initHook()
			go R.mainLoop()
		}
	}()

	defer func() { _ = windowsHook.MouseUnhook() }()
	defer func() { _ = windowsHook.KeyBoardUnhook() }()

	for {
		select {
		case kEvent := <-R.keyboardChan: //记录键盘操作
			R.recordKeyNode(&kEvent)
		case mEvent := <-R.mouseChan: //记录鼠标操作
			if !R.recordMouseTrack {
				if mEvent.Message == windowsHook.WM_MOUSEMOVE {
					R.lastMoveEven = &mEvent
					continue
				} else if R.lastMoveEven != nil {
					R.recordMouseNode(R.lastMoveEven)
				}
			}
			R.recordMouseNode(&mEvent)
		default:
			//别再犯那么傻的事情了，没有default会按顺序去尝试执行，然后卡住
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// ----------------------- Util -----------------------

// 初始化勾子
func (R *RecordServerT) initHook() {

	var err error
	if R.hotKeyChan, err = windowsHook.KeyBoardHook(nil); err != nil {
		R.tryPublishServerError(err)
		panic("记录勾子初始化失败. " + err.Error())
	}

	go func() {
		for event := range R.hotKeyChan {
			if !R.hotKeyDeal(event) {
				select {
				//尝试塞到键盘监听中，无效则丢弃
				case R.keyboardChan <- event:
				default:
				}
			}
		}
	}()
}

// 取消勾子
func (R *RecordServerT) handOutHook() {
	//鼠标直接取消勾子
	_ = windowsHook.MouseUnhook()
	//键盘将监听chan置空
	R.keyboardChan = nil
}

// 挂上勾子
func (R *RecordServerT) handUpHook() {
	//鼠标
	_ = windowsHook.MouseUnhook()

	var err error
	R.mouseChan, err = windowsHook.MouseHook(nil)
	R.tryPublishServerError(err)

	//键盘
	R.keyboardChan = make(chan windowsHook.KeyboardEvent, 3000)
}

// 热键处理
func (R *RecordServerT) hotKeyDeal(event windowsHook.KeyboardEvent) (isHotKey bool) {
	if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(event.VkCode)]; exist {
		isHotKey = true
		R.publishHotKeyDown(hotKey)
	}
	return
}

func (R *RecordServerT) recordKeyNode(event *windowsHook.KeyboardEvent) {

	R.notes = append(R.notes, noteT{
		NoteType: keyMouTool.TYPE_INPUT_KEYBOARD,
		KeyNote: &keyMouTool.KeyInputT{VK: keyMouTool.VKCode(event.VkCode),
			DwFlags: R.transKeyDwFlags(event.Message),
		},
		TimeGap: event.RecordTime - R.recordStartTime,
	})
	R.recordStartTime = event.RecordTime
}

func (R *RecordServerT) recordMouseNode(event *windowsHook.MouseEvent) {

	R.notes = append(R.notes, noteT{
		NoteType: keyMouTool.TYPE_INPUT_MOUSE,
		MouseNote: &keyMouTool.MouseInputT{X: event.X, Y: event.Y,
			DWFlags:   R.transMouseDwFlags(event.Message),
			MouseData: event.MouseData,
		},
		TimeGap: event.RecordTime - R.recordStartTime,
	})
	R.recordStartTime = event.RecordTime
}

// 记录到文件
func (R *RecordServerT) recordNoteToFile(name string, notes []noteT) {
	logTool.DebugAJ("record 开始记录文件：" + "名称:" + name + " 长度：" + strconv.Itoa(len(notes)))

	if name == "" {
		return
	}

	if len(notes) == 0 {
		return
	}

	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0772)
	if err != nil {
		R.tryPublishServerError(err)
		return
	}
	defer func() { _ = file.Close() }()

	js, err := json.Marshal(notes)
	if err != nil {
		R.tryPublishServerError(err)
		return
	}
	_, err = file.Write(js)
	if err != nil {
		R.tryPublishServerError(err)
		return
	}
}

func (R *RecordServerT) transMouseDwFlags(message windowsHook.Message) (dw keyMouTool.MouseInputDW) {

	return R.mouseDwMap[message]
}

func (R *RecordServerT) transKeyDwFlags(message windowsHook.Message) keyMouTool.KeyBoardInputDW {
	return R.keyDwMap[message]
}

// ----------------------- publish -----------------------

// 发布服务错误事件
func (R *RecordServerT) tryPublishServerError(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}

// 发布热键按下事件
func (R *RecordServerT) publishHotKeyDown(hotKey enum.HotKey) {
	err := eventCenter.Event.Publish(events.ServerHotKeyDown, events.ServerHotKeyDownData{
		HotKey: hotKey,
	})
	R.tryPublishServerError(err)
}
