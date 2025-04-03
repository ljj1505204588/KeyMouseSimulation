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
	"sort"
	"strconv"
	"strings"
	"time"
)

const FileExt = ".rpf"

type FileControlI interface {
	Save(name string, data []keyMouTool.NoteT)
	ReadFile(name string) (data keyMouTool.MulNote)

	Choose(name string) error
	Current() (name string, files []string)
}

var FileControl FileControlI

func init() {
	var f = fileControlT{}
	{
		f.getWindowRect()

		f.basePath, _ = os.Getwd()
		f.basePath = filepath.Join(f.basePath, "record")
		common.MustNil(os.MkdirAll(f.basePath, 0666))

		f.fileCheck()
		go f.scanFile()
	}
	FileControl = &f
}

type fileControlT struct {
	windowsX int // 电脑屏幕宽度
	windowsY int // 电脑屏幕长度

	current   string   // 当前文件
	fileNames []string // 所有文件

	basePath string // 基础路径

}

// Save 存储
func (f *fileControlT) Save(name string, data []keyMouTool.NoteT) {
	if name == "" || len(data) == 0 {
		return
	}

	if gene.In(f.fileNames, name) {
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

	f.fileNames = append(f.fileNames, name)
	f.current = name

	f.publishFileListChange()
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
	if !gene.In(f.fileNames, name) {
		return errors.New(language.ErrorSaveFileNameNilStr.ToString())
	}

	f.current = name
	return nil
}

// Current 当前选择文件
func (f *fileControlT) Current() (name string, files []string) {
	return f.current, f.fileNames
}

// 扫描文件
func (f *fileControlT) scanFile() {
	defer func() {
		_ = recover()
		go f.scanFile()
	}()

	var ticker = time.NewTicker(2 * time.Second)
	for {
		f.fileCheck()
		<-ticker.C
	}
}
func (f *fileControlT) fileCheck() {

	//遍历存储当前文件名字
	var infos []os.FileInfo
	var fs, _ = os.ReadDir(f.basePath)
	for _, per := range fs {
		if filepath.Ext(per.Name()) == FileExt {
			var info, _ = per.Info()
			infos = append(infos, info)
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime().Unix() < infos[j].ModTime().Unix()
	})
	// 按时间排序
	var names []string
	for _, info := range infos {
		name := filepath.Base(info.Name())
		names = append(names, strings.TrimSuffix(name, FileExt))
	}

	// 对比
	if len(names) != len(f.fileNames) || len(gene.Intersection(names, f.fileNames)) != len(names) {
		f.fileNames = names

		if !gene.In(names, f.current) {
			f.current = ""
			if len(names) > 0 {
				f.current = names[0]
			}
		}

		f.publishFileListChange()
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

func (f *fileControlT) publishFileListChange() {
	f.publishErr(eventCenter.Event.Publish(topic.FileListChange, &topic.FileListChangeData{
		ChooseFile: f.current,
		Files:      f.fileNames,
	}))
}

func (f *fileControlT) publishErr(err error) {
	if err == nil {
		return
	}

	_ = eventCenter.Event.Publish(topic.ServerError, &topic.ServerErrorData{
		ErrInfo: err.Error(),
	})
}
