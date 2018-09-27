package backend

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type AddAction actions.Action

func (this *AddAction) Run(params struct {
	Filename string
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

	backend := new(teaconfigs.ServerBackendConfig)
	backend.Address = params.Address

	server.Backends = append(server.Backends, backend)
	server.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success("保存成功")
}
