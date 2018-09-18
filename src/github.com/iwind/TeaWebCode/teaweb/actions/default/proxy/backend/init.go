package backend

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Module("").
			Prefix("/proxy/backend").
			Post("/add", new(AddAction)).
			Post("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Prefix("").
			End()
	})
}
