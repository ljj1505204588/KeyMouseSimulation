package component

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"errors"
	"sync"
)

func NewHK(name enum.HotKey, key string, exec func()) (HotKeyI, error) {
	var hk = &hotKeyT{name: name}
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

func GetHkByCode(code keyMouTool.VKCode) (HotKeyI, bool) {
	var hkI, ok = manage.hookM.Load(code)
	if !ok {
		return &hotKeyT{}, false
	}
	return hkI.(HotKeyI), true
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

	SetMethod(f func())
	ExecMethod()
}

type hotKeyT struct {
	name enum.HotKey

	key  string
	code keyMouTool.VKCode
	exec func()
}

func (h *hotKeyT) SetKey(key string) error {
	if h.key == key {
		return nil
	}

	var code = keyMouTool.VKCodeStringMap[key]
	if _, ok := manage.hookM.Load(code); ok {
		return errors.New(language.Center.Get(language.SetHotKeyErrMessageStr))
	}

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
	if h.exec != nil {
		h.exec()
	}
}
