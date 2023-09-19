package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
	"sync"
)

type HotKeyServerI interface {
	SetHotKey(k string, vks keyMouTool.VKCode)
}

type HotKeyServerT struct {
	once sync.Once
	l    sync.Mutex

	m map[keyMouTool.VKCode]string //热键信息
}

var t HotKeyServerT

func GetHotKeyServer() HotKeyServerI {
	t.once.Do(t.start)

	return &t
}

func (t *HotKeyServerT) start() {
	eventCenter.Event.Register(events.WindowsKeyBoardHook, t.hotKeyHandler)
}

// SetHotKey 设置热键
func (t *HotKeyServerT) SetHotKey(k string, vks keyMouTool.VKCode) {
	defer t.lockSelf()()

	t.m[vks] = k
}

// 热键勾子
func (t *HotKeyServerT) hotKeyHandler(data interface{}) (err error) {
	if !t.l.TryLock() {
		return
	}
	defer t.l.Unlock()

	var info = data.(events.WindowsKeyBoardHookData)
	if k, exist := t.m[keyMouTool.VKCode(info.Date.VkCode)]; exist {
		err = eventCenter.Event.Publish(events.ServerHotKeyDown, events.ServerHotKeyDownData{
			Key: &k,
		})
		t.tryPublishServerError(err)
	}
	return
}

// 发布错误事件
func (t *HotKeyServerT) tryPublishServerError(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}

func (t *HotKeyServerT) lockSelf() func() {
	t.l.Lock()
	return t.l.Unlock
}
