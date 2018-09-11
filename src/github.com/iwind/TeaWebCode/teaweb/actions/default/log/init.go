package log

import "github.com/iwind/TeaWebCode/teaplugin"

func init() {
	plugin := teaplugin.NewPlugin()

	widget := teaplugin.NewWidget()
	widget.Dashboard = true
	widget.Name = "QPS"
	widget.Title = "QPS/带宽"
	widget.URL = "/log/widget"
	plugin.AddWidget(widget)

	teaplugin.Register(plugin)
}
