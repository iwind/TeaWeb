package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaplugins"
	"github.com/iwind/TeaWebCode/teamongo"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	// 检查系统

	// 检查mongodb
	err := teamongo.Test()
	if err != nil {
		this.RedirectURL("/install/mongo")
		return
	}

	this.Data["teaMenu"] = "dashboard"

	widgets := teaplugins.DashboardWidgets()
	this.Data["widgets"] = widgets

	this.Show()
}
