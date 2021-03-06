package monitor

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/monitor").
			Get("", new(IndexAction)).
			EndAll()
	})
}
