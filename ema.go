package ematimeout

import (
	"math"
	"sync"
)

// 因为每个cgi有可能不一样，根据自身业务需求考虑是否要为每个接口create不同的ema计算
const (
	// 最低响应时间， 一般用平均响应时间替代
	Tavg = 55
	// 超时时间限制， 确保最坏的时候，所有请求能处理。正常时正确处理的成功率满足需求。
	Thwm = 300
	// 最大弹性时间
	Tmax = 500
	// 平滑指数，越小表示平均值越受最近值的影响，太大则对异常响应较慢
	N = 90
)

// Options define ema options
type Options struct {
	// 最低响应时间， 一般用平均响应时间替代 (ms)
	Tavg float64
	// 超时时间限制， 确保最坏的时候，所有请求能处理。正常时正确处理的成功率满足需求。 (ms)
	Thwm float64
	// 最大弹性时间 (ms)
	Tmax float64
	// 平滑指数，越小表示平无值越受最近值的影响，太大则对异常响应较慢
	N float64
}

// EMA instance of ema
type EMA struct {
	mu      sync.RWMutex
	options *Options
	ema     float64
	r       float64
}

// New new instance of EMA
func New() *EMA {
	defaultOptions := &Options{
		Tavg: Tavg,
		Thwm: Thwm,
		Tmax: Tmax,
		N:    N,
	}
	return &EMA{
		options: defaultOptions,
		ema:     0,
		r:       1 / (defaultOptions.N + 1),
	}
}

// NewFrom new instance of EMA
func NewFrom(options *Options) *EMA {
	return &EMA{
		options: options,
		ema:     0,
		r:       1 / (options.N + 1),
	}
}

// Update 更新最近延时时间, 计算当前ema
func (e *EMA) Update(x float64) float64 {
	e.mu.Lock()
	ema := x*e.r + e.ema*(1-e.r)
	e.ema = ema
	e.mu.Unlock()
	return ema
}

// Get 获取当前延时时间 (ms)
func (e *EMA) Get() float64 {
	tdto := 0.0 // 当前延时间间
	e.mu.RLock()
	ema := e.ema
	e.mu.RUnlock()
	if ema <= e.options.Tavg {
		// EMA <= Tavg， 这种情况EMA评估值比Tavg还要少，按耗时越少，给的延时越多的原则来计算，当前给的延时时间为Tmax。
		tdto = e.options.Tmax
	} else if ema >= e.options.Thwm {
		// EMA >= Thwm, 这种情况下，网络抖动比较严重，直接把延时限制在Thwm上，加快失败，方便进行柔性处理。
		tdto = e.options.Thwm
	} else {
		// EMA > Tavg, 这种情况EMA估值已经大于平均延时，此时启用动态延时计算
		// 平无表现相对限制的偏移比例，相当于弹性
		p := (e.options.Thwm - ema) / (e.options.Thwm - e.options.Tavg)
		tdto = e.options.Thwm + p*(e.options.Tmax-e.options.Thwm)
	}
	return math.Abs(tdto)
}
