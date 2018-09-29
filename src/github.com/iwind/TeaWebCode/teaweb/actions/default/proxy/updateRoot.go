package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type UpdateRootAction actions.Action

func (this *UpdateRootAction) Run(params struct {
	Filename string
	Root     string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Root = params.Root
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Success()
}
