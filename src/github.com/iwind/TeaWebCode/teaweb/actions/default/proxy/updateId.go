package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
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

	// 检查ID是否已经被使用
	for _, p := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
		if p.Filename == params.Filename {
			continue
		}
		if p.Id == params.Id {
			this.FailField("id", "此代理ID已经被使用，请换一个")
		}
	}

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Id = params.Id
	proxy.WriteToFile(Tea.ConfigFile(params.Filename))

	global.NotifyChange()

	this.Refresh().Success("保存成功")
}
