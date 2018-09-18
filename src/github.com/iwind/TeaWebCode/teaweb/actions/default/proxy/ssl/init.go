package ssl

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Module("").
			Prefix("/proxy/ssl").
			Post("/uploadCert", new(UploadCertAction)).
			Post("/uploadKey", new(UploadKeyAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			Prefix("").
			EndAll()
	})
}
