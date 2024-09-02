package param

import (
	"KeyMouseSimulation/common/param/paramcenter/file"
	"KeyMouseSimulation/common/param/paramoption"
	"encoding/json"
	"sync"
)

// ModuleI 配置模块接口·
type ModuleI interface {
	ParamName() (name string)
}

// CenterI 配置中心接口
// 当前实现：[本地文件]
type CenterI interface {
	Read(name string) (data any, exist bool)
	Write(name string, data any) (err error)
	UpdateRegister(name string, f func(name string, data any) error)
}

// ManagerI 管理中心接口
type ManagerI interface {
	Register(c ModuleI, options ...paramoption.Option)
	RegisterUpdate(name string, h func())
	SetOptions(name string, options ...paramoption.Option)
}

var Manager managerT

func init() {
	Manager.centers.Store("file", &paramcenter.F)
}

// ------------------------------------------------ 接口 ------------------------------------------------

type managerT struct {
	modules       sync.Map
	updateHandler sync.Map
	options       sync.Map

	centers sync.Map
}

func (m *managerT) SetOptions(name string, options ...paramoption.Option) {
	m.options.Store(name, options)
}

func (m *managerT) Register(mod ModuleI, options ...paramoption.Option) {
	var name = mod.ParamName()
	if _, have := m.modules.LoadOrStore(name, mod); have {
		panic("repeat register.")
	}
	m.options.Store(name, options)

	m.centers.Range(func(key, value any) bool {
		cen := value.(CenterI)
		if data, exist := cen.Read(name); exist {
			if err := m.parseModuleI(name, data); err != nil {
				return true
			}
		}

		if err := cen.Write(name, mod); err != nil {
			return true
		}
		cen.UpdateRegister(name, m.getUpdateHandler(cen))

		return true
	})
}

func (m *managerT) RegisterUpdate(name string, h func()) {

	var hs []func()
	var val, ok = m.updateHandler.Load(name)
	if ok {
		hs = val.([]func())
	}

	hs = append(hs, h)
	m.updateHandler.Store(name, hs)
}

func (m *managerT) getUpdateHandler(cen CenterI) func(name string, data any) (err error) {
	return func(name string, data any) error {
		if data, exist := cen.Read(name); exist {
			if err := m.parseModuleI(name, data); err != nil {
				return err
			}
		}

		if h, exist := m.updateHandler.Load(name); exist {
			for _, per := range h.([]func()) {
				per()
			}
		}
		return nil
	}
}

func (m *managerT) parseModuleI(name string, data any) (err error) {

	if mod, exist := m.modules.Load(name); exist {
		var bytes []byte
		if bytes, err = json.Marshal(data); err == nil {
			return json.Unmarshal(bytes, mod)
		}
	}
	return
}
