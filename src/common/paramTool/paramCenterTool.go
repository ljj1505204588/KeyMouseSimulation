package paramTool

import (
	"encoding/xml"
	"reflect"
	"unsafe"
)

var Center centerT

type centerT struct {
	param  paramT
	active map[ParamEnum]bool
}

// 创建
func (c *centerT) create() {
	c.active = defaultActive
	//初始化 && 接口校验
	var offset uintptr
	var value = reflect.ValueOf(c.param)
	for i := 0; i < value.NumField(); i++ {
		//初始化
		t := (**templateT)(unsafe.Pointer(uintptr(unsafe.Pointer(&c.param)) + offset))
		*t = (*templateT)(reflect.New(value.Field(i).Type().Elem()).UnsafePointer())
		offset += unsafe.Sizeof(value.Field(i).UnsafePointer())

		//接口校验
		if _, ok := value.Field(i).Interface().(templateT); !ok {
			panic(PARAM_TOOL_LOG_HEAD + value.Field(i).Type().Name() + " 未实现接口。")
		}
	}

	return
}

// Active 激活对应枚举
// 当基础模块使用param获取参数时候，若为默认未激活模块需要手动激活
func (c *centerT) Active(enum ParamEnum, writeBack bool) {

	c.active[enum] = true

	value := reflect.ValueOf(c.param)
	funParam := reflect.ValueOf(&c.param)

	//方法调用
	field := value.FieldByName(string(enum))
	field.MethodByName("EnvParamInit").Call([]reflect.Value{funParam})
	field.MethodByName("DefaultParamInit").Call([]reflect.Value{funParam})

	if writeBack {
		_ = WriteBackFile()
	}
	return
}

// UnActive 取消激活对应枚举
func (c *centerT) UnActive(enum ParamEnum) {
	delete(c.active, enum)
}

// GetParam 获取参数配置
func (c *centerT) GetParam() *paramT {
	return &c.param
}

// 初始化配置
func (c *centerT) initParam(text []byte) {
	_ = xml.Unmarshal(text, &c.param)

	for activeName := range defaultActive {
		c.Active(activeName, false)
	}

	if err := WriteBackFile(); err != nil {
		panic("配置回写失败。")
	}

	return
}

// 获取回写会文件夹的参数
func (c *centerT) getWriterBack() paramT {

	var back = paramT{}
	var backValue = reflect.ValueOf(back)
	var offSet uintptr
	for i := 0; i < backValue.NumField(); i++ {
		tag := backValue.Type().Field(i).Tag.Get("enum")

		if c.active[ParamEnum(tag)] {
			//通过unsafePoint 偏移去做
			t := (**templateT)(unsafe.Pointer(uintptr(unsafe.Pointer(&back)) + offSet))
			*t = *(**templateT)(unsafe.Pointer(uintptr(unsafe.Pointer(&c.param)) + offSet))
		}
		offSet += unsafe.Sizeof(backValue.Field(i).UnsafePointer())
	}

	return back
}
