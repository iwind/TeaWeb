package log

import (
	"github.com/iwind/TeaWebCode/teaplugins"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
)

func init() {
	plugin := teaplugins.NewPlugin()

	widget := teaplugins.NewWidget()
	widget.Dashboard = true
	widget.Name = "QPS"
	widget.Title = "QPS/带宽"
	widget.URL = "/log/widget"
	plugin.AddWidget(widget)

	teaplugins.Register(plugin)

	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			EndAll().
			Helper(new(helpers.UserMustAuth)).
			Prefix("/log").
			Get("", new(IndexAction)).
			Get("/get", new(GetAction)).
			GetPost("/widget", new(WidgetAction)).
			EndAll()
	})
}
