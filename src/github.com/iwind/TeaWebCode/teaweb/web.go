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
		Get("/", index.IndexAction{}).
		GetPost("/login", login.IndexAction{}).
		Get("/logout", logout.IndexAction{}).
		Get("/dashboard", dashboard.IndexAction{}).
		Get("/proxy", proxy.IndexAction{}).
		Get("/log", log.IndexAction{}).
		Get("/log/get", log.GetAction{}).
		Get("/settings", settings.IndexAction{}).

		Session(sessions.NewFileSessionManager(86400, "gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9")).
		Start()
}
