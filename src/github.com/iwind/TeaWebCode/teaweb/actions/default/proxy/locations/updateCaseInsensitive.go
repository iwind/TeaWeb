package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type UpdateCaseInsensitiveAction actions.Action

func (this *UpdateCaseInsensitiveAction) Run(params struct {
	Filename        string
	Index           int
	CaseInsensitive bool
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		location.SetPattern(location.PatternString(), location.PatternType(), params.CaseInsensitive, location.IsReverse())
		proxy.WriteToFilename(params.Filename)
	}

	global.NotifyChange()

	this.Success()
}
