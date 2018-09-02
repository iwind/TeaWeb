package teaweb

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/index"
	"github.com/iwind/TeaWebCode/teaproxy"
	"github.com/iwind/TeaGo/sessions"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/login"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/logout"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/dashboard"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/proxy"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/log"
	"github.com/iwind/TeaWebCode/teaweb/actions/@default/settings"
)

func Start() {
	// 启动代理
	teaproxy.Start()

	// 启动管理界面
	TeaGo.NewServer().
		AccessLog(false).
		Get("/", new(index.IndexAction)).
		GetPost("/login", new(login.IndexAction)).
		Get("/logout", new(logout.IndexAction)).
		Get("/dashboard", new(dashboard.IndexAction)).
		Get("/proxy", new(proxy.IndexAction)).
		Get("/log", new(log.IndexAction)).
		Get("/log/get", new(log.GetAction)).
		Get("/settings", new(settings.IndexAction)).

		Session(sessions.NewFileSessionManager(86400, "gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9")).
		Start()
}
