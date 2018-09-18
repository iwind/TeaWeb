package proxy

import (
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Module("").
			Prefix("/proxy").
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/detail", new(DetailAction)).
			Post("/updateDescription", new(UpdateDescriptionAction)).
			Post("/addName", new(AddNameAction)).
			Post("/updateName", new(UpdateNameAction)).
			Post("/deleteName", new(DeleteNameAction)).

			Post("/addListen", new(AddListenAction)).
			Post("/deleteListen", new(DeleteListenAction)).
			Post("/updateListen", new(UpdateListenAction)).

			Prefix("").
			End()
	})
}
