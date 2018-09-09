package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"time"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename string
	Auth     *helpers.UserMustAuth
}) {
	configFile := files.NewFile(Tea.ConfigFile(params.Filename))
	if !configFile.Exists() {
		this.Fail("要删除的配置文件不存在")
	}

	err := configFile.Delete()
	if err != nil {
		logs.Error(err)
		this.Fail("配置文件删除失败")
	}

	// 重启
	go func() {
		time.Sleep(1 * time.Second)
		teaproxy.Shutdown()
		teaproxy.Start()
	}()

	this.Refresh().Success()
}
