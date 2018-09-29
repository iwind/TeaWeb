package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type UpdateRootAction actions.Action

func (this *UpdateRootAction) Run(params struct {
	Filename string
	Index    int
	Root     string
	Must     *actions.Must
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		location.Root = params.Root
		proxy.WriteToFilename(params.Filename)

		global.NotifyChange()
	}

	this.Success()
}
