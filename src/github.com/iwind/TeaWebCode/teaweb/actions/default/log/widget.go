package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/tealogs"
)

type WidgetAction actions.Action

func (this *WidgetAction) Run(params struct{}) {
	this.Show()
}

func (this *WidgetAction) RunPost(params struct{}) {
	logger := tealogs.SharedLogger()
	this.Data["qps"] = logger.QPS()
	this.Data["inputBandwidth"] = logger.InputBandWidth()
	this.Data["outputBandwidth"] = logger.OutputBandWidth()

	this.Success()
}
