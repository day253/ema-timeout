package ematimeout

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestEMA_Normal(t *testing.T) {
	options := &Options{
		Tavg: 60,
		Thwm: 250,
		Tmax: 500,
		N:    90,
	}
	ema := NewFrom(options)
	for i := 0; i < 100; i++ {
		a := rand.Float64() * float64(200)
		e := ema.Update(a)
		t := ema.Get()
		fmt.Println("a: ", a, ", e: ", e, ", t: ", t)
	}
}

func TestEMA_Abnormal(t *testing.T) {
	options := &Options{
		Tavg: 60,
		Thwm: 250,
		Tmax: 500,
		N:    90,
	}
	ema := NewFrom(options)
	for i := 0; i < 100; i++ {
		a := rand.Float64()*float64(200) + float64(500)
		e := ema.Update(a)
		t := ema.Get()
		fmt.Println("a: ", a, ", e: ", e, ", t: ", t)
	}
}

func TestEMA_Burr(t *testing.T) {
	options := &Options{
		Tavg: 55,
		Thwm: 300,
		Tmax: 500,
		N:    90,
	}
	ema := NewFrom(options)
	for i := 0; i < 100; i++ {
		a := rand.Float64()*float64(200) + float64(500)
		e := ema.Update(a)
		t := ema.Get()
		fmt.Println("a: ", a, ", e: ", e, ", t: ", t)
	}
}
