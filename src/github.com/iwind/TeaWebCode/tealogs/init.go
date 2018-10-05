package tealogs

import (
	"github.com/iwind/TeaWebCode/teaplugins"
	"github.com/iwind/TeaWebCode/teacharts"
	"github.com/iwind/TeaGo/timers"
	"time"
)

func init() {
	plugin := teaplugins.NewPlugin()
	createWidget(plugin)
	teaplugins.Register(plugin)
}

func createWidget(plugin *teaplugins.Plugin) {
	widget := teaplugins.NewWidget()
	widget.Dashboard = true
	widget.Group = teaplugins.WidgetGroupRealTime
	widget.Name = "日志"

	chart := teacharts.NewTable()
	chart.Name = "即时日志"
	accessLogs := SharedLogger().ReadNewLogs("", 10)
	for _, accessLog := range accessLogs {
		chart.AddRow("<em>" + accessLog.TimeLocal + "-" + accessLog.Host + "</em><br/> \"" + accessLog.Request + "\"")
	}
	widget.AddChart(chart)

	plugin.AddWidget(widget)

	// 定时刷新
	timers.Loop(3*time.Second, func(looper *timers.Looper) {
		accessLogs := SharedLogger().ReadNewLogs("", 5)
		chart.ClearRows()
		for _, accessLog := range accessLogs {
			chart.AddRow("<em>" + accessLog.TimeLocal + "</em><br/> \"" + accessLog.Request + "\"")
		}
	})
}
