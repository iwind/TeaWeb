package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaWebCode/teaproxy"
	"github.com/iwind/TeaGo/logs"
)

type DeleteListenAction actions.Action

func (this *DeleteListenAction) Run(params struct {
	Filename string
	Index    int
	Auth     *helpers.UserMustAuth
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(proxy.Listen) {
		list := lists.NewList(proxy.Listen)
		list.Remove(params.Index)
		proxy.Listen = list.Slice.([]string)
	}

	logs.Println(proxy.Listen)
	proxy.WriteToFilename(params.Filename)

	// 重启服务
	teaproxy.Restart()

	this.Refresh().Success()
}
