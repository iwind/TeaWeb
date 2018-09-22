package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
)

type UpdateIdAction actions.Action

func (this *UpdateIdAction) Run(params struct {
	Filename string
	Id       string
	Must     *actions.Must
}) {
	params.Must.
		Field("id", params.Id).
		Require("代理ID不能为空")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Id = params.Id
	proxy.WriteToFile(Tea.ConfigFile(params.Filename))

	this.Refresh().Success("保存成功")
}
