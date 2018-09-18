package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type AddNameAction actions.Action

func (this *AddNameAction) Run(params struct {
	Filename string
	Name     string
	Must     *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入域名")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.AddName(params.Name)
	proxy.WriteToFilename(params.Filename)

	// 重启服务
	teaproxy.Restart()

	this.Refresh().Success("保存成功")
}
