package login

import "github.com/iwind/TeaGo"

func init()  {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Module("").
			Prefix("/login").
			GetPost("", new(IndexAction)).
			Prefix("").
			End()
	})
}