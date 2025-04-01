package rp_file

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/common/gene"
	"KeyMouseSimulation/common/windowsApi"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/topic"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const FileExt = ".rpf"

type FileControlI interface {
	Save(name string, data []keyMouTool.NoteT)
	ReadFile(name string) (data keyMouTool.MulNote)

	Choose(name string) error
	Current() (name string)
}

var FileControl FileControlI

func init() {
	var f = fileControlT{}
	{
		f.getWindowRect()

		f.basePath, _ = os.Getwd()
		f.basePath = filepath.Join(f.basePath, "record")
		common.MustNil(os.MkdirAll(f.basePath, 0666))

		go f.scanFile()
	}
	FileControl = &f
}

type fileControlT struct {
	windowsX int // 电脑屏幕宽度
	windowsY int // 电脑屏幕长度

	current  string   // 当前文件
	fileName []string // 所有文件

	basePath string // 基础路径

}

// Save 存储
func (f *fileControlT) Save(name string, data []keyMouTool.NoteT) {
	if name == "" || len(data) == 0 {
		return
	}

	if gene.In(f.fileName, name) {
		name += "_" + strconv.Itoa(int(time.Now().Unix()))
	}

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
	f.current = name

	if err = eventCenter.Event.Publish(topic.FileListChange, &topic.FileListChangeData{
		ChooseFile: name,
		Files:      f.fileName,
	}); err != nil {
		f.publishErr(err)
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

	data.AdaptWindow(f.windowsX, f.windowsY) // 适应屏幕

	return
}

// Choose 文件选择
func (f *fileControlT) Choose(name string) error {
	if !gene.In(f.fileName, name) {
		return errors.New(language.ErrorSaveFileNameNilStr.ToString())
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
	defer func() {
		_ = recover()
		go f.scanFile()
	}()

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
			f.fileName = names

			if !gene.In(names, f.current) {
				f.current = ""
				if len(names) > 0 {
					f.current = names[0]
				}
			}

			f.publishErr(eventCenter.Event.Publish(topic.FileListChange, &topic.FileListChangeData{
				ChooseFile: f.current,
				Files:      f.fileName,
			}))

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
	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: err.Error(),
	})
}
