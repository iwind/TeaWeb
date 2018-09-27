package ssl

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/global"
)

type OnAction actions.Action

func (this *OnAction) Run(params struct {
	Filename string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if server.SSL == nil {
		ssl := new(teaconfigs.SSLConfig)
		ssl.On = true
		server.SSL = ssl
	} else {
		server.SSL.On = true
	}

	server.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Success()
}
