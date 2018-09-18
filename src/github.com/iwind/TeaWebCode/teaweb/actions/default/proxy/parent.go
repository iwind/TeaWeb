package proxy

import (
	"github.com/iwind/TeaWebCode/teaweb/actions/utils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/lists"
)

type ParentAction struct {
	utils.ParentAction
}

func (this *ParentAction) Before() {
	this.Data["teaMenu"] = "proxy"
	this.Data["teaTabbar"] = []maps.Map{
		{
			"name":    "已有代理",
			"subName": "",
			"url":     this.URL("/proxy"),
			"active":  lists.Contains([]string{"proxy.IndexAction", "proxy.DetailAction", "proxy.UpdateAction"}, this.Spec.ClassName),
		},
		{
			"name":    "添加新代理",
			"subName": "",
			"url":     this.URL("/proxy/add"),
			"active":  this.Spec.ClassName == "proxy.AddAction",
		},
	}
}
