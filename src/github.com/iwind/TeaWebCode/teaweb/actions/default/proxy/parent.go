package proxy

import (
	"github.com/iwind/TeaWebCode/teaweb/actions/utils"
	"github.com/iwind/TeaGo/maps"
)

type ParentAction struct {
	utils.ParentAction
}

func (this *ParentAction) Before() {
	this.Data["teaMenu"] = "proxy"
	this.Data["teaTabbar"] = []maps.Map{
		{
			"name":    "已有服务",
			"subName": "",
			"url":     this.URL("/proxy"),
			"active":  this.Spec.ClassName == "proxy.IndexAction" || this.Spec.ClassName == "proxy.UpdateAction",
		},
		{
			"name":    "添加新服务",
			"subName": "",
			"url":     this.URL("/proxy/add"),
			"active":  this.Spec.ClassName == "proxy.AddAction",
		},
	}
}
