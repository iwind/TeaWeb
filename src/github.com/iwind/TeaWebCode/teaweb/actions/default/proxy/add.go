package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

// 添加新的服务
type AddAction actions.Action

func (this *AddAction) Run(params struct {
}) {
	this.Show()
}

func (this *AddAction) RunPost(params struct {
	Id          string
	Description string
	Must        *actions.Must
}) {
	// ID是否已存在
	if len(params.Id) > 0 {
		ids := maps.Map{}
		for _, config := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
			ids[config.Id] = true
		}
		if ids.Has(params.Id) {
			this.FailField("id", "此代理ID已经被使用，请换一个")
		}
	} else {
		params.Id = stringutil.Rand(8)
	}

	// 描述
	if len(params.Description) == 0 {
		params.Description = "新服务"
	}

	server := teaconfigs.NewServerConfig()
	server.Http = true
	server.Id = params.Id
	server.Description = params.Description
	server.Charset = "utf-8"

	filename := stringutil.Rand(16) + ".proxy.conf"
	configPath := Tea.ConfigFile(filename)
	err := server.WriteToFile(configPath)
	if err != nil {
		this.Fail(err.Error())
	}

	global.NotifyChange()

	this.Next("/proxy/detail", map[string]interface{}{
		"filename": filename,
	}, "").Success("添加成功，现在去配置详细信息")
}
