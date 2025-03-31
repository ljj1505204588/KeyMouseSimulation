package windowsInput

type MouseInputT struct {
	Type uint32
	Mi   MouseInputMiT
}
type MouseInputMiT struct {
	X           int32
	Y           int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type KeyboardInputT struct {
	Type uint32
	Ki   KeyboardInputKiT
}
type KeyboardInputKiT struct {
	WVk         int16
	WScan       int16
	DwFlags     uint32
	Time        int32
	DwExtraInfo uintptr
	unused      [8]byte
}

type HardWareInputT struct {
	Type uint32
	Hi   HardWareInputHiT
}
type HardWareInputHiT struct {
	UMsg    uint32
	WParamL uint16
	WParamH uint16
	unused  [16]byte
}
