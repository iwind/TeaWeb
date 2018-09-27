package rewrite

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type AddAction actions.Action

func (this *AddAction) Run(params struct {
	Filename   string
	Index      int
	Pattern    string
	Replace    string
	ProxyId    string
	TargetType string
	Must       *actions.Must
}) {
	//@TODO proxyId 支持一个Host

	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入匹配规则").

		Field("targetType", params.TargetType).
		In([]string{"url", "proxy"}, "目标类型错误")

	if params.TargetType == "proxy" {
		params.Must.
			Field("proxyId", params.ProxyId).
			Require("请选择目标代理")
	}

	params.Must.
		Field("replace", params.Replace).
		Require("请输入目标URL")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if len(params.Replace) == 0 || params.Replace[0] != '/' {
		params.Replace = "/" + params.Replace
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		rewriteRule := teaconfigs.NewRewriteRule()
		rewriteRule.On = true
		rewriteRule.Pattern = params.Pattern
		if params.TargetType == "url" {
			rewriteRule.Replace = params.Replace
		} else {
			rewriteRule.Replace = "proxy://" + params.ProxyId + params.Replace
		}
		location.Rewrite = append(location.Rewrite, rewriteRule)
	}

	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success("添加成功")
}
