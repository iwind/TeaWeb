package log

import (
	"github.com/iwind/TeaWebCode/teaplugin"
	"github.com/iwind/TeaGo"
)

func init() {
	plugin := teaplugin.NewPlugin()

	widget := teaplugin.NewWidget()
	widget.Dashboard = true
	widget.Name = "QPS"
	widget.Title = "QPS/带宽"
	widget.URL = "/log/widget"
	plugin.AddWidget(widget)

	teaplugin.Register(plugin)

	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Module("").
			Prefix("/log").
			Get("", new(IndexAction)).
			Get("/get", new(GetAction)).
			GetPost("/widget", new(WidgetAction)).
			Prefix("").
			End()
	})
}
