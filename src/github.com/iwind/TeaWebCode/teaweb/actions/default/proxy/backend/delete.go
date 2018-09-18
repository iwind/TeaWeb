package backend

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Auth     *helpers.UserMustAuth
	Filename string
	Index    int
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(server.Backends) {
		backends := lists.NewList(server.Backends)
		backends.Remove(params.Index)

		server.Backends = backends.Slice.([]*teaconfigs.ServerBackendConfig)
	}

	server.WriteToFilename(params.Filename)
	teaproxy.Restart()

	this.Refresh().Success()
}
