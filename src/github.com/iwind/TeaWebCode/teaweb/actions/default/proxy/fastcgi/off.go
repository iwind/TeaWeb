package fastcgi

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type OffAction actions.Action

func (this *OffAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	fastcgi := location.Fastcgi
	if fastcgi == nil {
		this.Fail("没有fastcgi配置，请刷新后重试")
	}

	fastcgi.On = false
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Success()
}
