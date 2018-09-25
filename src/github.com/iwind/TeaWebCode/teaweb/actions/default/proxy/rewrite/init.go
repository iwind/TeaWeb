package rewrite

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/rewrite").
			Helper(new(helpers.UserMustAuth)).
			Helper(new(proxy.Helper)).
			Post("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Post("/update", new(UpdateAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			EndAll()
	})
}
