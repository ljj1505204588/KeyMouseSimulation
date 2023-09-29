package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
	"sync"
)

type HotKeyServerT struct {
	once sync.Once
	l    sync.Mutex

	m map[keyMouTool.VKCode]string // 热键信息

}

var t HotKeyServerT

func init() {
	t.once.Do(t.start)
}

func (t *HotKeyServerT) start() {
	t.m = make(map[keyMouTool.VKCode]string)

	eventCenter.Event.Register(events.WindowsKeyBoardHook, t.hotKeyHandler)
	eventCenter.Event.Register(events.SetHotKey, t.setHotKey)
}

// SetHotKey 设置热键
func (t *HotKeyServerT) setHotKey(data interface{}) (err error) {
	defer t.lockSelf()()
	var info = data.(events.SetHotKeyData)

	if key,ok := keyMouTool.VKCodeStringMap[info.Key];ok {
		t.m[key] = info.Key
	}

	return
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
			Key: k,
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
