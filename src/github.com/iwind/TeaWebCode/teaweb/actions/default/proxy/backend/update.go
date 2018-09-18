package backend

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type UpdateAction actions.Action

func (this *UpdateAction) Run(params struct {
	Filename string
	Index    int
	Address  string
	Must     *actions.Must
}) {
	params.Must.
		Field("address", params.Address).
		Require("请输入后端服务器地址")

	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(server.Backends) {
		server.Backends[params.Index].Address = params.Address
	}

	server.WriteToFilename(params.Filename)

	teaproxy.Restart()

	this.Refresh().Success("保存成功")
}
