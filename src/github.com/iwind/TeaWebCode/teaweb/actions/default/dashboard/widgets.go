package dashboard

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaplugins"
)

type WidgetsAction actions.Action

func (this *WidgetsAction) Run() {
	groups := teaplugins.DashboardGroups()
	for _, group := range groups {
		group.Reload()
	}
	this.Data["widgetGroups"] = groups

	this.Success()
}
