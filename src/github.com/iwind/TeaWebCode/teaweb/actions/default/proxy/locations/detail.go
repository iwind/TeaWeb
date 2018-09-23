package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
)

type DetailAction actions.Action

func (this *DetailAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index < 0 || params.Index >= len(proxy.Locations) {
		this.Fail("找不到要修改的路径配置")
	}

	this.Data["filename"] = params.Filename
	this.Data["locationIndex"] = params.Index
	this.Data["location"] = proxy.Locations[params.Index]
	this.Data["proxy"] = proxy

	this.Show()
}
