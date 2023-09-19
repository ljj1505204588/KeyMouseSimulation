package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/logTool"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type FileControlI interface {
	Save(name string, data []noteT)
	ReadFile(name string) (data []noteT)
}

var fileControl FileControlT

func GetFileControl() FileControlI {
	fileControl.once.Do(fileControl.scanFile)
	return &fileControl
}

type FileControlT struct {
	once sync.Once

	windowsX int // 电脑屏幕宽度
	windowsY int // 电脑屏幕长度

	fileName []string
}

// Save 存储
func (f *FileControlT) Save(name string, data []noteT) {
	// 记录到文件
	logTool.DebugAJ("record 开始记录文件：" + "名称:" + name + " 长度：" + strconv.Itoa(len(data)))

	if name == "" || len(data) == 0 {
		return
	}

	// 打开文件
	var file, err = os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0772)
	if err != nil {
		f.publishErr(err)
		return
	}
	defer func() { _ = file.Close() }()

	// 序列化
	js, err := json.Marshal(data)
	if err != nil {
		f.publishErr(err)
		return
	}

	// 存储
	if _, err = file.Write(js); err != nil {
		f.publishErr(err)
		return
	}
}

// ReadFile 读取文件
func (f *FileControlT) ReadFile(name string) (data []noteT) {
	var dealErr = func(err error) []noteT {
		f.publishErr(err)
		return nil
	}

	// ----------------------------

	file, err := os.OpenFile(name, os.O_RDONLY, 0772)
	if err != nil {
		return dealErr(err)
	}
	defer func() { _ = file.Close() }()

	b, err := io.ReadAll(file)
	if err != nil {
		return dealErr(err)
	}

	nodes := make(mulNote, 100)
	if err = json.Unmarshal(b, &nodes); err != nil {
		return dealErr(err)
	}

	nodes.adaptWindow(f.windowsX, f.windowsY)

	if err == nil {
		logTool.DebugAJ("playback 加载文件成功：" + "名称:" + name + " 长度：" + strconv.Itoa(len(nodes)))
	}

	return nodes
}

// 扫描文件
func (f *FileControlT) scanFile() {
	go func() {
		defer f.scanFile()
		for {
			var names []string
			//遍历存储当前文件名字
			if fs, err := os.ReadDir("./"); err == nil {
				for _, per := range fs {
					if filepath.Ext(per.Name()) == FileExt {
						name := filepath.Base(per.Name())
						name = name[:len(name)-len(FileExt)]
						names = append(names, name)
					}
				}
			}

			//对比
			if newFile := GenTool.Exclude(names, f.fileName); len(newFile) != 0 || len(names) != len(f.fileName) {
				f.fileName = names
				_ = eventCenter.Event.Publish(events.FileScanNewFile, events.FileScanNewFileData{
					NewFile:  newFile,
					FileList: names,
				})
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// 获取windows窗口大小
func (f *FileControlT) getWindowRect() {
	f.windowsX, f.windowsY = 1920, 1080

	x, _, err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics, windowsApi.SM_CXSCREEN)
	if err != nil {
		return
	}
	y, _, err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics, windowsApi.SM_CYSCREEN)
	if err != nil {
		return
	}

	f.windowsX, f.windowsY = int(x), int(y)
}

// ---------------------------------- 发布事件 ----------------------------------

func (f *FileControlT) publishErr(err error) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: err.Error(),
	})
}

// --------------------------------------------------------- 记录 ---------------------------------------------------------

type noteT struct {
	NoteType  keyMouTool.InputType
	KeyNote   *keyMouTool.KeyInputT
	MouseNote *keyMouTool.MouseInputT
	TimeGap   int64 //Nanosecond
	timeGap   float64
}

type mulNote []noteT

// 添加记录
func (m *mulNote) appendMouseNote(startTime int64, event *windowsHook.MouseEvent) {
	var dw, exist = mouseDwMap[event.Message]
	if exist {
		var note = noteT{
			NoteType: keyMouTool.TYPE_INPUT_MOUSE,
			MouseNote: &keyMouTool.MouseInputT{X: event.X, Y: event.Y,
				DWFlags:   dw,
				MouseData: event.MouseData,
			},
			TimeGap: max(startTime-event.RecordTime, 0),
		}

		*m = append(*m, note)
	}
}

// 添加记录
func (m *mulNote) appendKeyBoardNote(startTime int64, event *windowsHook.KeyboardEvent) {
	var dw, exist = keyDwMap[event.Message]
	if exist {
		var note = noteT{
			NoteType: keyMouTool.TYPE_INPUT_KEYBOARD,
			KeyNote: &keyMouTool.KeyInputT{VK: keyMouTool.VKCode(event.VkCode),
				DwFlags: dw,
			},
			TimeGap: max(startTime-event.RecordTime, 0),
		}

		*m = append(*m, note)
	}
}

// 适应窗口大小
func (m *mulNote) adaptWindow(x, y int) {

	for nodePos := range *m {
		(*m)[nodePos].timeGap = float64((*m)[nodePos].TimeGap)
		if (*m)[nodePos].NoteType == keyMouTool.TYPE_INPUT_MOUSE {
			(*m)[nodePos].MouseNote.X = (*m)[nodePos].MouseNote.X * 65535 / int32(x)
			(*m)[nodePos].MouseNote.Y = (*m)[nodePos].MouseNote.Y * 65535 / int32(y)
		}
	}
}

var mouseDwMap = map[windowsHook.Message]keyMouTool.MouseInputDW{
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
var keyDwMap = map[windowsHook.Message]keyMouTool.KeyBoardInputDW{
	windowsHook.WM_KEYDOWN: keyMouTool.DW_KEYEVENTF_KEYDown,
	windowsHook.WM_KEYUP:   keyMouTool.DW_KEYEVENTF_KEYUP,
}
