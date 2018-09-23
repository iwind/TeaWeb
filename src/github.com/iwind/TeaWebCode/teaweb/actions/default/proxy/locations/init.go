package locations

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Prefix("/proxy/locations").
			Helper(new(helpers.UserMustAuth)).
			Helper(new(proxy.Helper)).
			Get("", new(IndexAction)).
			Post("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Get("/detail", new(DetailAction)).
			EndAll()
	})
}
