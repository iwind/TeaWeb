package proxy

import (
	"github.com/iwind/TeaWebCode/teaconfigs"
)

type DetailAction struct {
	ParentAction
}

func (this *DetailAction) Run(params struct {
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
