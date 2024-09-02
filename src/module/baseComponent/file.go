package component

import (
	eventCenter "KeyMouseSimulation/common/Event"
	gene "KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/commonTool"
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/share/events"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"encoding/json"
	"errors"
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
	Save(name string, data []keyMouTool.NoteT)
	ReadFile(name string) (data keyMouTool.MulNote)

	FileChange(exec func(names, newNames []string))
	Choose(name string) error
	Current() (name string)
}

var FileControl FileControlI

func init() {
	var f = fileControlT{}
	f.once.Do(func() {
		f.getWindowRect()

		f.basePath, _ = os.Getwd()
		f.basePath = filepath.Join(f.basePath, "record")
		commonTool.MustNil(os.MkdirAll(f.basePath, 0666))

		go f.scanFile()
	})
	FileControl = &f
}

type fileControlT struct {
	once sync.Once

	changeExec func(names, newNames []string)

	windowsX int // 电脑屏幕宽度
	windowsY int // 电脑屏幕长度

	basePath string
	fileName []string

	current string
}

// Save 存储
func (f *fileControlT) Save(name string, data []keyMouTool.NoteT) {
	if name == "" || len(data) == 0 {
		return
	}
	logTool.DebugAJ(fmt.Sprintf("record 开始记录文件[%s],长度[%d]", name, len(data)))

	// 打开文件
	var filePath = filepath.Join(f.basePath, name+FileExt)
	var file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0772)
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

	f.fileName = append(f.fileName, name)
	if f.changeExec != nil {
		f.changeExec(f.fileName, []string{name})
	}
}

// ReadFile 读取文件
func (f *fileControlT) ReadFile(name string) (data keyMouTool.MulNote) {
	var dealErr = func(err error) []keyMouTool.NoteT {
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

	data.AdaptWindow(f.windowsX, f.windowsY)

	logTool.DebugAJ(fmt.Sprintf("playback 加载文件[%s]成功,长度[%d]", name, len(data)))
	return
}

// FileChange 文件变动执行回调
func (f *fileControlT) FileChange(exec func(names, newNames []string)) {
	f.changeExec = exec
}

// Choose 文件选择
func (f *fileControlT) Choose(name string) error {
	if !gene.Contain(f.fileName, name) {
		return errors.New(language.Center.Get(language.ErrorSaveFileNameNilStr))
	}

	f.current = name
	return nil
}

// Current 当前选择文件
func (f *fileControlT) Current() (name string) {
	return f.current
}

// 扫描文件
func (f *fileControlT) scanFile() {
	defer func() { go f.scanFile() }()
	for {
		var names []string
		//遍历存储当前文件名字
		if fs, err := os.ReadDir(f.basePath); err == nil {
			for _, per := range fs {
				if filepath.Ext(per.Name()) == FileExt {
					name := filepath.Base(per.Name())
					names = append(names, strings.TrimSuffix(name, FileExt))
				}
			}
		}

		// 对比
		if len(names) != len(f.fileName) || len(gene.Intersection(names, f.fileName)) != len(names) {
			if f.changeExec != nil {
				newFile := gene.Exclude(names, f.fileName)
				f.fileName = names
				f.changeExec(names, newFile)
			}
		}

		time.Sleep(2 * time.Second)
	}
}

// 获取windows窗口大小
func (f *fileControlT) getWindowRect() {
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

func (f *fileControlT) publishErr(err error) {
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: err.Error(),
	})
}
