package backend

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaWebCode/teaproxy"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename string
	Index    int
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 0 && params.Index < len(server.Backends) {
		server.Backends = lists.Remove(server.Backends, params.Index).([]*teaconfigs.ServerBackendConfig)
	}

	server.WriteToFilename(params.Filename)
	teaproxy.Restart()

	this.Refresh().Success()
}
