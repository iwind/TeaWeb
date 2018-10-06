package teaservices

import (
	"github.com/iwind/TeaWebCode/teaplugins"
)

type ProbeInterface interface {
	Run() // 检测
}

type Probe struct {
	isInitialized bool
	Plugin        *teaplugins.Plugin
	IsRunning     bool
}

func (this *Probe) InitOnce(f func()) {
	if this.isInitialized {
		return
	}
	this.isInitialized = true
	this.Plugin = teaplugins.NewPlugin()
	teaplugins.Register(this.Plugin)
	f()
}

func (this *Probe) Finish() {
	this.IsRunning = false
}
