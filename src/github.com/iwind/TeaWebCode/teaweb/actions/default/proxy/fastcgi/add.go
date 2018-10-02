package fastcgi

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type AddAction actions.Action

func (this *AddAction) Run(params struct {
	Filename    string
	Index       int
	On          bool
	Pass        string
	ReadTimeout int
	Params      string
	Must        *actions.Must
}) {
	params.Must.
		Field("filename", params.Filename).
		Require("请输入配置文件名").
		Field("pass", params.Pass).
		Require("请输入Fastcgi地址")

	paramsMap := map[string]string{}
	err := ffjson.Unmarshal([]byte(params.Params), &paramsMap)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	location.AddFastcgi(&teaconfigs.FastcgiConfig{
		On:          params.On,
		Pass:        params.Pass,
		ReadTimeout: fmt.Sprintf("%ds", params.ReadTimeout),
		Params:      paramsMap,
	})
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success()
}
