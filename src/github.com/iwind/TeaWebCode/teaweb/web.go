package teaweb

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/index"
	"github.com/iwind/TeaGo/sessions"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/login"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/logout"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/dashboard"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/proxy"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/log"
	"github.com/iwind/TeaWebCode/teaweb/actions/default/settings"
	"time"
	"github.com/iwind/TeaWebCode/teaproxy"
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
		GetPost("/login", new(login.IndexAction)).
		Get("/logout", new(logout.IndexAction)).
		Get("/dashboard", new(dashboard.IndexAction)).

		Get("/proxy", new(proxy.IndexAction)).
		GetPost("/proxy/add", new(proxy.AddAction)).
		Post("/proxy/delete", new(proxy.DeleteAction)).
		GetPost("/proxy/update", new(proxy.UpdateAction)).

		Get("/log", new(log.IndexAction)).
		Get("/log/get", new(log.GetAction)).
		GetPost("/log/widget", new(log.WidgetAction)).
		Get("/settings", new(settings.IndexAction)).

		Session(sessions.NewFileSessionManager(86400, "gSeDQJJ67tAVdnguDAQdGmnDVrjFd2I9")).
		Start()
}
