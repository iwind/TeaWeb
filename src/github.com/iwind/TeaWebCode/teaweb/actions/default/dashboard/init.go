package dashboard

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Prefix("/dashboard").
			Get("", new(IndexAction)).
			Get("/widgets", new(WidgetsAction)).
			EndAll()
	})
}
