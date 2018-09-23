package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
)

type UpdateReverseAction actions.Action

func (this *UpdateReverseAction) Run(params struct {
	Filename string
	Index    int
	Reverse  bool
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		location.SetPattern(location.PatternString(), location.PatternType(), location.IsCaseInsensitive(), params.Reverse)
		proxy.WriteToFilename(params.Filename)
	}

	this.Success()
}
