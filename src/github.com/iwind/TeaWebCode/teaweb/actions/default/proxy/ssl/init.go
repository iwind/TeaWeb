package ssl

import (
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Module("").
			Prefix("/proxy/ssl").
			Post("/uploadCert", new(UploadCertAction)).
			Post("/uploadKey", new(UploadKeyAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			Prefix("").
			End()
	})
}
