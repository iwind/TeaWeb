package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type AddListenAction actions.Action

func (this *AddListenAction) Run(params struct {
	Filename string
	Listen   string
	Must     *actions.Must
}) {
	params.Must.
		Field("listen", params.Listen).
		Require("请输入监听地址")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.AddListen(params.Listen)
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success("保存成功")
}
