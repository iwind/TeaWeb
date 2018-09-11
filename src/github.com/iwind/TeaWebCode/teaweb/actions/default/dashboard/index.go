package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaWebCode/teaplugin"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	this.Data["teaMenu"] = "dashboard"

	widgets := teaplugin.DashboardWidgets()
	this.Data["widgets"] = widgets

	this.Show()
}
