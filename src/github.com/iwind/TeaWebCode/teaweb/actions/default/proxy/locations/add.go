package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type AddAction actions.Action

func (this *AddAction) Run(params struct {
	Filename        string
	Pattern         string
	TypeId          int
	Reverse         bool
	CaseInsensitive bool
	Must            *actions.Must
}) {
	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入规则").
		Field("typeId", params.TypeId).
		Gt(0, "请选择类型").
		In([]int{
			teaconfigs.LocationPatternTypePrefix,
			teaconfigs.LocationPatternTypeExact,
			teaconfigs.LocationPatternTypeRegexp,
		}, "选择的类型错误")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := teaconfigs.NewLocationConfig()
	location.On = true
	location.SetPattern(params.Pattern, params.TypeId, params.CaseInsensitive, params.Reverse)

	proxy.Locations = append(proxy.Locations, location)
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Next("/proxy/locations/detail", map[string]interface{}{
		"filename": params.Filename,
		"index":    len(proxy.Locations) - 1,
	})

	this.Success("添加成功")
}
