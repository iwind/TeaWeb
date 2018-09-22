package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"time"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename string
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

	// @TODO 删除对应的certificate file和certificate key file

	// 重启
	go func() {
		time.Sleep(1 * time.Second)
		teaproxy.Restart()
	}()

	this.Refresh().Success()
}
