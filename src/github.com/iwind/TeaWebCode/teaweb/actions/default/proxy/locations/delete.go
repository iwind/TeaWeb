package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Locations) {
		proxy.Locations = lists.Remove(proxy.Locations, params.Index).([]*teaconfigs.LocationConfig)
	}

	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success()
}
