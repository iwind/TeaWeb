package teaweb

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/sessions"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/logout"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/dashboard"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/settings"
	"time"
	"github.com/iwind/TeaWebCode/teaproxy"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/install"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/index"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/login"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/proxy"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/ssl"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/backend"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/locations"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/proxy/rewrite"
	_ "github.com/iwind/TeaWebCode/teaweb/actions/default/log"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

func Start() {
	// 启动代理
	go func() {
		time.Sleep(1 * time.Second)
		teaproxy.Start()
	}()

	// 启动管理界面
	TeaGo.NewServer().
		AccessLog(false).
		Get("/", new(index.IndexAction)).
		Get("/logout", new(logout.IndexAction)).

		Helper(new(helpers.UserMustAuth)).
		Get("/dashboard", new(dashboard.IndexAction)).
		Get("/settings", new(settings.IndexAction)).
		GetPost("/install/mongo", new(install.MongoAction)).
		EndAll().

		Session(sessions.NewFileSessionManager(
			86400,
			"gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9",
		)).
		Start()
}
