package ssl

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type OffAction actions.Action

func (this *OffAction) Run(params struct {
	Filename string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if server.SSL != nil {
		server.SSL.On = false
	}

	server.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Success()
}
