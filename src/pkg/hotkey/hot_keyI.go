//go:build windows
// +build windows

package hk

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func (c *centerT) DefaultShowSign() (res map[enum.HotKey]string) {
	return map[enum.HotKey]string{
		enum.HotKeyRecord:   "F7",
		enum.HotKeyPlayBack: "F8",
		enum.HotKeyPause:    "F9",
		enum.HotKeyStop:     "F10",
	}
}

// Show 显示文字
func (c *centerT) Show() (res map[enum.HotKey]string) {
	res = make(map[enum.HotKey]string)
	for _, key := range enum.TotalHotkey() {
		switch key {
		case enum.HotKeyRecord:
			res[key] = language.RecordStr.ToString()
		case enum.HotKeyPlayBack:
			res[key] = language.PlayBackStr.ToString()
		case enum.HotKeyPause:
			res[key] = language.PauseStr.ToString()
		case enum.HotKeyStop:
			res[key] = language.StopStr.ToString()
		}
	}

	return
}

// ShowSign 显示热键标识
func (c *centerT) ShowSign() (res map[enum.HotKey]string) {
	defer common.LockSelf(&Center.lock)()

	res = make(map[enum.HotKey]string)
	// 遍历所有值
	for _, key := range enum.TotalHotkey() {
		if hk, ok := Center.hotKeyMap[key]; ok {
			res[key] = hk.show()
		}
	}

	return
}

//  --------------------------------- 热键中心 ---------------------------------

var Center = &centerT{
	hotKeyMap: make(map[enum.HotKey]hotKeyI),
}

func init() {
	_ = Center.hotKeySetHandler(&topic.HotKeySetData{
		Set: map[enum.HotKey]keyMouTool.VKCode{
			enum.HotKeyRecord:   keyMouTool.VKCodeStringMap["F7"],
			enum.HotKeyPlayBack: keyMouTool.VKCodeStringMap["F8"],
			enum.HotKeyPause:    keyMouTool.VKCodeStringMap["F9"],
			enum.HotKeyStop:     keyMouTool.VKCodeStringMap["F10"],
		},
	})
	eventCenter.Event.Register(topic.HotKeySet, Center.hotKeySetHandler)

	// 监听系统退出信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		Center.cleanup()
		fmt.Println("热键回调已删除.")
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
	var dataValue = data.(*topic.HotKeySetData)

	for HotKey, code := range dataValue.Set {
		// 取出维护的
		var hk, ok = c.hotKeyMap[HotKey]
		if !ok {
			hk = &hotKeyT{}
			c.hotKeyMap[HotKey] = hk
		}

		// 设置
		hk.set(HotKey, 0, code)
	}

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
