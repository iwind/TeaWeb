package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaplugins"
)

type WidgetsAction actions.Action

func (this *WidgetsAction) Run() {
	this.Data["widgetGroups"] = teaplugins.DashboardGroups()

	this.Success()
}
