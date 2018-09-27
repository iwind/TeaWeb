package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaproxy"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type RestartAction actions.Action

func (this *RestartAction) Run(params struct{}) {
	teaproxy.Restart()

	global.FinishChange()

	this.Success()
}
