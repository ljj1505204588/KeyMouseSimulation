package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/logTool"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
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
	Save(name string ,data []noteT)
	ReadFile(name string )(data []noteT)
}

var fileControl FileControlT

func GetFile()FileControlI{
 	fileControl.once.Do(fileControl.scanFile)
	return &fileControl
}

type FileControlT struct {
	once sync.Once

	windowsX int // 电脑屏幕宽度
	windowsY int // 电脑屏幕长度

	fileName []string
}

func (f *FileControlT)Save(name string ,data []noteT){
	// 记录到文件
	logTool.DebugAJ("record 开始记录文件：" + "名称:" + name + " 长度：" + strconv.Itoa(len(data)))

	if name == "" || len(data) == 0{
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

func (f *FileControlT)ReadFile(name string )(data []noteT){
	var dealErr = func(err error) []noteT{
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
	if err = json.Unmarshal(b, &nodes); err != nil{
		return dealErr(err)
	}

	nodes.adaptWindow(f.windowsX,f.windowsY)

	if err == nil {
		logTool.DebugAJ("playback 加载文件成功：" + "名称:" + name + " 长度：" + strconv.Itoa(len(nodes)))
	}

	return nodes
}

func (f *FileControlT)scanFile(){
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
			if newFile :=GenTool.Exclude(names, f.fileName) ; len(newFile) != 0 || len(names) != len(f.fileName){
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

func (f *FileControlT)publishErr(err error){
	_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
		ErrInfo: err.Error(),
	})
}