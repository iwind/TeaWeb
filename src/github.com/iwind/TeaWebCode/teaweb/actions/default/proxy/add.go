package proxy

import (
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaGo/actions"
	"strings"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/utils/string"
	"time"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type AddAction struct {
	ParentAction
}

func (this *AddAction) Run(params struct {
	Auth *helpers.UserMustAuth
}) {
	this.Show()
}

func (this *AddAction) RunPost(params struct {
	Auth           *helpers.UserMustAuth
	Name           []string
	ListenAddress  []string
	ListenPort     []int
	BackendAddress []string
	BackendPort    []int
	Must           *actions.Must
}) {
	if len(params.Name) == 0 {
		this.Fail("域名不能为空")
	}

	for index, name := range params.Name {
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			this.Fail("域名不能为空")
		}
		params.Name[index] = name
	}

	for index, address := range params.ListenAddress {
		address = strings.TrimSpace(address)
		if len(address) == 0 {
			this.Fail("访问地址不能为空")
		}
		params.ListenAddress[index] = address
	}

	for index, port := range params.ListenPort {
		if port <= 0 || port >= 65535 {
			this.Fail("访问地址端口错误")
		}
		params.ListenPort[index] = port
	}

	for index, address := range params.BackendAddress {
		address = strings.TrimSpace(address)
		if len(address) == 0 {
			this.Fail("后端地址不能为空")
		}
		params.BackendAddress[index] = address
	}

	for index, port := range params.BackendPort {
		if port <= 0 || port >= 65535 {
			this.Fail("后端地址端口错误")
		}
		params.BackendPort[index] = port
	}

	// 保存
	server := teaconfigs.NewServerConfig()
	server.AddName(params.Name ...)
	for index, address := range params.ListenAddress {
		if index > len(params.ListenPort)-1 {
			continue
		}

		server.AddListen(fmt.Sprintf("%s:%d", address, params.ListenPort[index]))
	}
	for index, address := range params.BackendAddress {
		if index > len(params.BackendPort)-1 {
			continue
		}

		backend := &teaconfigs.ServerBackendConfig{
			Address: fmt.Sprintf("%s:%d", address, params.BackendPort[index]),
		}
		server.AddBackend(backend)
	}

	err := server.WriteToFile(Tea.ConfigFile(stringutil.Rand(16) + ".proxy.conf"))
	if err != nil {
		this.Fail("配置文件写入失败")
	}

	// 重启
	go func() {
		time.Sleep(1 * time.Second)
		teaproxy.Restart()
	}()

	this.Next("/proxy", nil, "").Success("服务保存成功")
}
