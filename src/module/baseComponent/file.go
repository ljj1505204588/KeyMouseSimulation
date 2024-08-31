package component

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/logTool"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/share/events"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const FileExt = ".rpf"

type FileControlI interface {
	Save(name string, data []NoteT)
	ReadFile(name string) (data MulNote)
}

var FileControl FileControlI

func init() {
	var f = FileControlT{}
	f.once.Do(func() {
		f.getWindowRect()
		f.basePath, _ = os.Getwd()
		go f.scanFile()
	})
	FileControl = &f
}

type FileControlT struct {
	once sync.Once

	windowsX int // 电脑屏幕宽度
	windowsY int // 电脑屏幕长度

	basePath string
	fileName []string
}

// Save 存储
func (f *FileControlT) Save(name string, data []NoteT) {
	if name == "" || len(data) == 0 {
		return
	}
	logTool.DebugAJ(fmt.Sprintf("record 开始记录文件[%s],长度[%d]", name, len(data)))

	// 打开文件
	name = filepath.Join(f.basePath, name+FileExt)
	var file, err = os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0772)
	if err != nil {
		f.publishErr(err)
		return
	}
	defer func() { _ = file.Close() }()

	// 存储
	js, _ := json.Marshal(data)
	if _, err = file.Write(js); err != nil {
		f.publishErr(err)
	}
}

// ReadFile 读取文件
func (f *FileControlT) ReadFile(name string) (data MulNote) {
	var dealErr = func(err error) []NoteT {
		f.publishErr(err)
		return nil
	}

	// ----------------------------

	file, err := os.OpenFile(filepath.Join(f.basePath, name+FileExt), os.O_RDONLY, 0772)
	if err != nil {
		return dealErr(err)
	}
	defer func() { _ = file.Close() }()

	b, err := io.ReadAll(file)
	if err != nil {
		return dealErr(err)
	}

	if err = json.Unmarshal(b, &data); err != nil {
		return dealErr(err)
	}

	data.adaptWindow(f.windowsX, f.windowsY)

	logTool.DebugAJ(fmt.Sprintf("playback 加载文件[%s]成功,长度[%d]", name, len(data)))
	return
}

// 扫描文件
func (f *FileControlT) scanFile() {
	defer func() { go f.scanFile() }()
	for {
		var names []string
		//遍历存储当前文件名字
		if fs, err := os.ReadDir("./"); err == nil {
			for _, per := range fs {
				if filepath.Ext(per.Name()) == FileExt {
					name := filepath.Base(per.Name())
					names = append(names, strings.TrimSuffix(name, FileExt))
				}
			}
		}

		//对比
		if newFile := GenTool.Exclude(names, f.fileName); len(newFile) != 0 {
			if err := eventCenter.Event.Publish(events.FileScanNewFile, events.FileScanNewFileData{
				NewFile:  newFile,
				FileList: names,
			}); err == nil {
				f.fileName = names
			}
		}
		time.Sleep(2 * time.Second)
	}
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
