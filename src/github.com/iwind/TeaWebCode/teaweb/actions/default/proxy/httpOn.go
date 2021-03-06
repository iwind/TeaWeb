package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaWebCode/teaconfigs"
)

type HttpOnAction actions.Action

func (this *HttpOnAction) Run(params struct {
	Filename string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Http = true
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Success()
}
