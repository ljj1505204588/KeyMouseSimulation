package paramcenter

import (
	fsnotify "gopkg.in/fsnotify.v1"
	"sync"
	"time"
)

var F fileCenterT

func init() {
	F.fileDeal.init()
	F.updateHandler = make(map[string]func(name string, data any) error)

	var err error
	if F.watcher, err = fsnotify.NewWatcher(); err != nil {
		panic(err)
	}

	go F.WatchWorker()
}

type fileCenterT struct {
	fileDeal      fileDealT // 文件处理
	watcher       *fsnotify.Watcher
	updateHandler map[string]func(name string, data any) error

	sync.Mutex
}

func (f *fileCenterT) Read(name string) (data any, exist bool) {
	defer f.lockSelf()()

	if cur, err := f.fileDeal.Read(); err == nil {
		for _, per := range cur {
			if per.Name == name {
				return per.Data, true
			}
		}
	}

	return
}
func (f *fileCenterT) Write(name string, data any) error {
	defer f.lockSelf()()

	var cur, _ = f.fileDeal.Read()

	var add = true
	for i, per := range cur {
		if per.Name == name {
			add = false
			cur[i].Data = data
		}
	}

	if add {
		cur = append(cur, modData{
			Name: name,
			Data: data,
		})
	}

	return f.fileDeal.WriteBack(cur)
}

func (f *fileCenterT) UpdateRegister(name string, h func(name string, data any) error) {
	defer f.lockSelf()()

	f.updateHandler[name] = h
}

func (f *fileCenterT) WatchWorker() {
	defer func() {
		if rec := recover(); rec != nil {
			go f.WatchWorker()
		}
	}()

	err := F.watcher.Add(F.fileDeal.path)
	for err != nil {
		err = F.watcher.Add(F.fileDeal.path)
		time.Sleep(2 * time.Second)
	}

	for {
		select {
		case event := <-f.watcher.Events:
			switch event.Op {
			case fsnotify.Write:
				f.update()
			}
		case _ = <-f.watcher.Errors:

		}
	}
}
func (f *fileCenterT) update() {

	for _, handler := range f.updateHandler {
		if cur, err := f.fileDeal.Read(); err == nil {
			for _, per := range cur {
				_ = handler(per.Name, per.Data)
			}
		}
	}
}

// -------------------------------- util --------------------------------

func (f *fileCenterT) lockSelf() func() {
	f.Lock()
	return f.Unlock
}
