package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/utils/string"
)

type AddAction struct {
	ParentAction
}

func (this *AddAction) Run(params struct {
}) {
	this.Show()
}

func (this *AddAction) RunPost(params struct {
	Description string
	Must        *actions.Must
}) {
	params.Must.
		Field("description", params.Description).
		Require("代理说明不能为空")

	server := teaconfigs.NewServerConfig()
	server.Description = params.Description

	filename := stringutil.Rand(16) + ".proxy.conf"
	configPath := Tea.ConfigFile(filename)
	err := server.WriteToFile(configPath)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Next("/proxy/detail", map[string]interface{}{
		"filename": filename,
	}, "").Success()
}
