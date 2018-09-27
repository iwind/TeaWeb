package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type StatusAction actions.Action

func (this *StatusAction) Run() {
	this.Data["changed"] = global.ProxyIsChanged()
	this.Success()
}
