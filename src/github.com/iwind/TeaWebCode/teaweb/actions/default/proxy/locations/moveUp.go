package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type MoveUpAction actions.Action

func (this *MoveUpAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 1 && params.Index < len(proxy.Locations) {
		prev := proxy.Locations[params.Index-1]
		current := proxy.Locations[params.Index]
		proxy.Locations[params.Index-1] = current
		proxy.Locations[params.Index] = prev
	}

	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success()
}
