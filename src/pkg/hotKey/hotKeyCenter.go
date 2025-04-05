//go:build windows
// +build windows

package hk

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/common/elegantExit"
	"KeyMouseSimulation/common/gene"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/pkg/language"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"fmt"
	"sync"
)

func (c *centerT) IsHotKey(vk keyMouTool.VKCode) bool {
	defer common.LockSelf(&c.lock)()
	return c.vkCodeMap[vk]
}

func (c *centerT) DefaultShowSign() (res map[enum.HotKey]string) {
	return map[enum.HotKey]string{
		enum.HotKeyRecord:    "F7",
		enum.HotKeyPlayBack:  "F8",
		enum.HotKeyPause:     "F9",
		enum.HotKeyStop:      "F10",
		enum.HotKeySpeedUp:   ">",
		enum.HotKeySpeedDown: "<",
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
		case enum.HotKeySpeedUp:
			res[key] = language.SpeedUpStr.ToString()
		case enum.HotKeySpeedDown:
			res[key] = language.SpeedDownStr.ToString()
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
	vkCodeMap: make(map[keyMouTool.VKCode]bool),
}

func init() {
	var (
		recordConf    = conf.RecordHotKeyConf.GetValue()
		playBackConf  = conf.PlayBackHotKeyConf.GetValue()
		pauseConf     = conf.PauseHotKeyConf.GetValue()
		stopConf      = conf.StopHotKeyConf.GetValue()
		speedUpConf   = conf.SpeedUpHotKeyConf.GetValue()
		speedDownConf = conf.SpeedDownHotKeyConf.GetValue()
	)
	_ = Center.hotKeySetHandler(&topic.HotKeySetData{
		Set: map[enum.HotKey]string{
			enum.HotKeyRecord:    gene.Choose(recordConf != "", recordConf, "F7"),
			enum.HotKeyPlayBack:  gene.Choose(playBackConf != "", playBackConf, "F8"),
			enum.HotKeyPause:     gene.Choose(pauseConf != "", pauseConf, "F9"),
			enum.HotKeyStop:      gene.Choose(stopConf != "", stopConf, "F10"),
			enum.HotKeySpeedUp:   gene.Choose(speedUpConf != "", speedUpConf, ">"),
			enum.HotKeySpeedDown: gene.Choose(speedDownConf != "", speedDownConf, "<"),
		},
	})
	eventCenter.Event.Register(topic.HotKeySet, Center.hotKeySetHandler)

	elegantExit.AddElegantExit(func() {
		Center.cleanup()
		fmt.Println("热键回调已删除.")
	})
}

type centerT struct {
	hotKeyMap map[enum.HotKey]hotKeyI
	vkCodeMap map[keyMouTool.VKCode]bool
	lock      sync.Mutex
}

// 设置热键事件回调函数
func (c *centerT) hotKeySetHandler(data interface{}) (err error) {
	defer common.LockSelf(&c.lock)()
	var dataValue = data.(*topic.HotKeySetData)

	for HotKey, sign := range dataValue.Set {
		// 取出维护的
		var hk, ok = c.hotKeyMap[HotKey]
		if !ok {
			hk = &hotKeyT{}
			c.hotKeyMap[HotKey] = hk
		}

		// 设置
		hk.set(HotKey, 0, keyMouTool.VKCodeStringMap[sign])

		// todo 想想
		switch HotKey {
		case enum.HotKeyRecord:
			conf.RecordHotKeyConf.SetValue(sign)
		case enum.HotKeyPlayBack:
			conf.PlayBackHotKeyConf.SetValue(sign)
		case enum.HotKeyPause:
			conf.PauseHotKeyConf.SetValue(sign)
		case enum.HotKeyStop:
			conf.StopHotKeyConf.SetValue(sign)
		case enum.HotKeySpeedUp:
			conf.SpeedUpHotKeyConf.SetValue(sign)
		case enum.HotKeySpeedDown:
			conf.SpeedDownHotKeyConf.SetValue(sign)
		}

	}

	// 更新vkcodeMap
	c.vkCodeMap = make(map[keyMouTool.VKCode]bool)
	for _, hk := range c.hotKeyMap {
		c.vkCodeMap[hk.vkCode()] = true
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
