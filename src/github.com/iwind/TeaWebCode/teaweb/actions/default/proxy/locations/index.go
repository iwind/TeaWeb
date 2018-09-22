package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Filename string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["filename"] = params.Filename
	this.Data["proxy"] = proxy

	this.Show()
}
