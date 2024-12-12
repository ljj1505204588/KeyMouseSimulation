package component

import (
	eventCenter "KeyMouseSimulation/common/Event"
	gene "KeyMouseSimulation/common/GenTool"
	"KeyMouseSimulation/common/share/enum"
	"KeyMouseSimulation/common/share/events"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"errors"
	"sync"
	"time"
)

func NewHK(name enum.HotKey, key string, exec func()) (HotKeyI, error) {
	var hk = &hotKeyT{name: name, defKey: key}
	manage.nameM.Store(name, hk)

	hk.exec = exec
	return hk, hk.SetKey(key)
}
func GetHk(name enum.HotKey) (HotKeyI, bool) {
	var hkI, ok = manage.nameM.Load(name)
	if !ok {
		return &hotKeyT{}, false
	}

	return hkI.(HotKeyI), true
}

func GetAllHk() (res map[enum.HotKey]HotKeyI) {
	res = make(map[enum.HotKey]HotKeyI)
	manage.nameM.Range(func(key, value any) bool {
		res[key.(enum.HotKey)] = value.(HotKeyI)
		return true
	})
	return
}

func GetHkByCode(code keyMouTool.VKCode) (HotKeyI, bool) {
	var hkI, ok = manage.hookM.Load(code)
	if !ok {
		return &hotKeyT{}, false
	}
	return hkI.(HotKeyI), true
}

func MulSetKey(mul map[HotKeyI]string) (err error) {
	var saveKeys []string
	var current = make(map[string]HotKeyI)
	for hk, key := range mul {
		saveKeys = append(saveKeys, key)
		current[hk.Key()] = hk
	}
	if len(gene.RemoveDuplicate(saveKeys)) != len(saveKeys) {
		return errors.New(Center.Get(SetHotKeyErrMessageStr))
	}

	for hk, key := range mul {
		// 如果当前目标key的hk也是需要更改的，先删除。
		if currentHk, ok := current[key]; ok && currentHk != hk {
			var code = keyMouTool.VKCodeStringMap[key]
			manage.hookM.Delete(code)
			if err = hk.SetKey(key); err != nil {
				manage.hookM.Store(key, currentHk)
				return
			}
		} else {
			if err = hk.SetKey(key); err != nil {
				return
			}
		}
	}

	return nil
}

func init() {
	eventCenter.Event.Register(events.WindowsKeyBoardHook, hotKeyHandler)
}

var manage = manageT{}

type manageT struct {
	nameM sync.Map
	hookM sync.Map
}

// 热键勾子
func hotKeyHandler(data interface{}) (err error) {
	var info = data.(events.WindowsKeyBoardHookData)
	if k, exist := manage.hookM.Load(keyMouTool.VKCode(info.Date.VkCode)); exist {
		go k.(HotKeyI).ExecMethod()
	}
	return
}

// ------------------------------ manage ------------------------------

type HotKeyI interface {
	SetKey(key string) error
	Key() string
	DefaultKey() string

	SetMethod(f func())
	ExecMethod()
}

type hotKeyT struct {
	name enum.HotKey

	defKey   string            // 默认Key
	key      string            // 当前Key
	code     keyMouTool.VKCode // 对应windowsCode
	exec     func()            // 执行方法
	lastExec time.Time         // 上次执行时间
}

func (h *hotKeyT) DefaultKey() string {
	return h.defKey
}

func (h *hotKeyT) SetKey(key string) error {
	if h.key == key {
		return nil
	}

	// 冲突判断
	var code = keyMouTool.VKCodeStringMap[key]
	if _, ok := manage.hookM.Load(code); ok {
		return errors.New(Center.Get(SetHotKeyErrMessageStr))
	}

	// 旧值删除
	if h.key != "" {
		manage.hookM.Delete(h.code)
	}

	h.key = key
	h.code = keyMouTool.VKCodeStringMap[key]
	manage.hookM.Store(h.code, h)
	return nil
}

func (h *hotKeyT) Key() string {
	return h.key
}

func (h *hotKeyT) SetMethod(f func()) {
	h.exec = f
}
func (h *hotKeyT) ExecMethod() {
	if h.exec == nil {
		return
	}

	// 执行间隔
	var execTime = time.Now()
	if time.Since(h.lastExec) > 300*time.Millisecond {
		h.exec()
		h.lastExec = execTime
	}
}
