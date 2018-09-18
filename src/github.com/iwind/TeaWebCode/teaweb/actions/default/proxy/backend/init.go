package backend

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Module("").
			Prefix("/proxy/backend").
			Post("/add", new(AddAction)).
			Post("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Prefix("").
			EndAll()
	})
}
