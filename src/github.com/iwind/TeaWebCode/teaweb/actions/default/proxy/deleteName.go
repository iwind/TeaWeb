package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type DeleteNameAction actions.Action

func (this *DeleteNameAction) Run(params struct {
	Filename string
	Index    int
	Auth     *helpers.UserMustAuth
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Name) {
		list := lists.NewList(proxy.Name)
		list.Remove(params.Index)
		proxy.Name = list.Slice.([]string)
	}

	proxy.WriteToFilename(params.Filename)

	// 重启服务
	teaproxy.Restart()

	this.Refresh().Success()
}
