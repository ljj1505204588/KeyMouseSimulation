package component

import "KeyMouseSimulation/common/param"

func init() {
	// 记录
	var recordParam = recordParamT{}
	recordParam.defValue()                                                          // 默认配置
	RecordConfig = &recordConfigT{recordParamT: recordParam}                        // 外部访问接口
	param.Manager.Register(&recordParam)                                            // 注册配置
	param.Manager.RegisterUpdate(recordParam.ParamName(), RecordConfig.paramUpdate) // 注册配置更新监听

	// 回放
	var playbackParam = playbackParamT{}
	playbackParam.defValue()                                                            // 默认配置
	PlaybackConfig = &playbackConfigT{playbackParamT: playbackParam}                    // 外部访问接口
	param.Manager.Register(&playbackParam)                                              // 注册配置
	param.Manager.RegisterUpdate(playbackParam.ParamName(), PlaybackConfig.paramUpdate) // 注册配置更新监听
}

var RecordConfig recordConfigI
var PlaybackConfig playbackConfigI

// ----------------------------------- 记录配置 -----------------------------------

type recordConfigI interface {
	paramUpdate()
	SetMouseTrackChange(exec bool, method func(record bool)) // 鼠标路径记录变动回调
	SetMouseTrack(record bool)                               // 设置是否记录鼠标路径
}

type recordConfigT struct {
	mouseTrackMethod []func(record bool)
	recordParamT
}

func (r *recordConfigT) SetMouseTrack(record bool) {
	r.MouseTrack = record
	for _, method := range r.mouseTrackMethod {
		method(record)
	}
}
func (r *recordConfigT) SetMouseTrackChange(exec bool, method func(record bool)) {
	r.mouseTrackMethod = append(r.mouseTrackMethod, method)
	if exec {
		method(r.MouseTrack)
	}
}
func (r *recordConfigT) paramUpdate() {
	for _, method := range r.mouseTrackMethod {
		method(r.MouseTrack)
	}
}

type recordParamT struct {
	MouseTrack bool // 鼠标路径
}

func (r *recordParamT) ParamName() (name string) {
	return "record_param"
}

func (r *recordParamT) defValue() {
	r.MouseTrack = true
}

// ----------------------------------- 回放配置 -----------------------------------
type playbackConfigI interface {
	paramUpdate()
	SetSpeed(speed float64)                               // 速度设置
	SetSpeedChange(exec bool, method func(speed float64)) // 速度变动回调

	SetPlaybackTimes(times int64)                               // 回放次数设置
	SetPlaybackTimesChange(exec bool, method func(times int64)) // 回放次数变动回调

	SetPlaybackRemainTimes(times int64)                               // 回放剩余次数
	SetPlaybackRemainTimesChange(exec bool, method func(times int64)) // 回放剩余次数回调
}

type playbackConfigT struct {
	speedMethod []func(speed float64)

	playbackTimesMethod       []func(times int64)
	playbackRemainTimesMethod []func(remainTimes int64)

	playbackParamT
}

func (p *playbackConfigT) SetSpeed(speed float64) {
	p.Speed = speed
	for _, method := range p.speedMethod {
		method(speed)
	}
}
func (p *playbackConfigT) SetSpeedChange(exec bool, method func(speed float64)) {
	p.speedMethod = append(p.speedMethod, method)
	if exec {
		method(p.Speed)
	}
}

func (p *playbackConfigT) SetPlaybackTimes(times int64) {
	p.PlaybackTimes = times
	for _, method := range p.playbackTimesMethod {
		method(times)
	}
}
func (p *playbackConfigT) SetPlaybackTimesChange(exec bool, method func(times int64)) {
	p.playbackTimesMethod = append(p.playbackTimesMethod, method)
	if exec {
		method(p.PlaybackTimes)
	}
}

func (p *playbackConfigT) SetPlaybackRemainTimes(times int64) {
	p.PlaybackRemainTimes = times
	for _, method := range p.playbackRemainTimesMethod {
		method(times)
	}
}
func (p *playbackConfigT) SetPlaybackRemainTimesChange(exec bool, method func(times int64)) {
	p.playbackRemainTimesMethod = append(p.playbackRemainTimesMethod, method)
	if exec {
		method(p.PlaybackRemainTimes)
	}
}

func (p *playbackConfigT) paramUpdate() {
	for _, method := range p.speedMethod {
		method(p.Speed)
	}

	for _, method := range p.playbackTimesMethod {
		method(p.PlaybackTimes)
	}
	for _, method := range p.playbackRemainTimesMethod {
		method(p.PlaybackRemainTimes)
	}
}

type playbackParamT struct {
	Speed float64 // 回放速度

	PlaybackTimes       int64 // 回放次数
	PlaybackRemainTimes int64 // 回放剩余次数
}

func (r *playbackParamT) ParamName() (name string) {
	return "playback_param"
}
func (r *playbackParamT) defValue() {
	r.Speed = 1

	r.PlaybackTimes = 1
	r.PlaybackRemainTimes = 0
}
