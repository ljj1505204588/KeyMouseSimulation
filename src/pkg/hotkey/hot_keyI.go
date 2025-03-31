//go:build windows
// +build windows

package hk

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/event_topic"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Show 显示文字
func Show(key enum.HotKey) string {
	switch key {
	case enum.HotKeyRecord:
		return language.RecordStr.ToString()
	case enum.HotKeyPlayBack:
		return language.PlayBackStr.ToString()
	case enum.HotKeyPause:
		return language.PauseStr.ToString()
	case enum.HotKeyStop:
		return language.StopStr.ToString()
	}
	return ""
}

// ShowSign 显示热键标识
func ShowSign(key enum.HotKey) string {
	defer common.LockSelf(&center.lock)()

	if hk, ok := center.hotKeyMap[key]; ok {
		return hk.show()
	}

	return ""
}

//  --------------------------------- 热键中心 ---------------------------------

var center = &centerT{
	hotKeyMap: make(map[enum.HotKey]hotKeyI),
}

func init() {
	eventCenter.Event.Register(event_topic.HotKeySet, center.hotKeySetHandler)

	// 监听系统退出信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		center.cleanup()
		os.Exit(0)
	}()
}

type centerT struct {
	hotKeyMap map[enum.HotKey]hotKeyI
	lock      sync.Mutex
}

// 设置热键事件回调函数
func (c *centerT) hotKeySetHandler(data interface{}) (err error) {
	defer common.LockSelf(&c.lock)()
	var dataValue = data.(event_topic.HotKeySetData)

	// 取出维护的
	var hk, ok = c.hotKeyMap[dataValue.HotKey]
	if !ok {
		hk = &hotKeyT{}
		c.hotKeyMap[dataValue.HotKey] = hk
	}

	// 设置
	hk.set(dataValue.HotKey, 0, dataValue.KeyBoardCodes)
	return
}

// 清理所有热键
func (c *centerT) cleanup() {
	defer common.LockSelf(&c.lock)()
	for _, hk := range c.hotKeyMap {
		if h, ok := hk.(*hotKeyT); ok {
			h.cleanup()
		}
	}
	c.hotKeyMap = make(map[enum.HotKey]hotKeyI)
}

type hotKeyI interface {
	set(key enum.HotKey, mod, code keyMouTool.VKCode) // 设置热键
	show() string                                     // 显示热键标识
}
