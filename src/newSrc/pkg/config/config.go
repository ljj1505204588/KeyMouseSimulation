package conf

import (
	"KeyMouseSimulation/common/commonTool"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Config 配置管理接口
type Config interface {
	Get(key ConfigKeyI) (any, bool)
	Set(key ConfigKeyI, value any) error
	Save() error
	Load() error
	OnChange(key ConfigKeyI, handler func(any))
}

// configuration 配置实现
type configuration struct {
	lock           sync.RWMutex
	path           string
	changeHandlers map[string][]func(any)
	values         map[string]any
}

// New 创建配置实例
func New(options ...ConfigOption) Config {
	c := &configuration{
		changeHandlers: make(map[string][]func(any)),
		values:         make(map[string]any),
		path:           "config.json",
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// Get 获取配置值
func (c *configuration) Get(key ConfigKeyI) (any, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if value, exists := c.values[key.getKey()]; exists {
		return value, true
	}
	defaultValue, _ := key.getDefault()
	return defaultValue, false
}

// Set 设置配置值
func (c *configuration) Set(key ConfigKeyI, value any) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := key.validate(value); err != nil {
		return err
	}

	c.values[key.getKey()] = value
	c.notifyHandlers(key.getKey(), value)
	return nil
}

// OnChange 注册配置变更处理器
func (c *configuration) OnChange(confKey ConfigKeyI, handler func(any)) {
	defer commonTool.RLockSelf(&(c.lock))()

	var key = confKey.getKey()
	c.changeHandlers[key] = append(c.changeHandlers[key], handler)
}

// Save 保存配置
func (c *configuration) Save() error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(c.values, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, bytes, 0644)
}

// Load 加载配置
func (c *configuration) Load() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	bytes, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return c.Save()
		}
		return err
	}

	newValues := make(map[string]any)
	if err := json.Unmarshal(bytes, &newValues); err != nil {
		return err
	}

	// 验证所有值
	for _, key := range allKeys() {
		if value, exists := newValues[key.getKey()]; exists {
			if err := key.validate(value); err != nil {
				return err
			}
		}
	}

	// 更新值并触发通知
	for k, v := range newValues {
		c.values[k] = v
		c.notifyHandlers(k, v)
	}

	return nil
}

// notifyHandlers 通知配置变更
func (c *configuration) notifyHandlers(key string, value any) {
	defer commonTool.RRLockSelf(&c.lock)()

	handlers := c.changeHandlers[key]

	for _, handler := range handlers {
		go handler(value)
	}
}
