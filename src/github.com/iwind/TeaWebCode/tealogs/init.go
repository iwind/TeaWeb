package tealogs

import (
	"github.com/iwind/TeaWebCode/teaplugins"
	"github.com/iwind/TeaWebCode/teacharts"
	"github.com/iwind/TeaGo/logs"
	"math"
)

func init() {
	logs.Println("register log plugin")

	plugin := teaplugins.NewPlugin()
	createWidget(plugin)
	teaplugins.Register(plugin)
}

func createWidget(plugin *teaplugins.Plugin) {
	widget := teaplugins.NewWidget()
	widget.Dashboard = true
	widget.Group = teaplugins.WidgetGroupRealTime
	widget.Name = "即时访问及流量"
	widget.MoreURL = "/log"

	// 带宽
	outputBandWidthChart := teacharts.NewGaugeChart()
	outputBandWidthChart.Name = "Web出口带宽"
	outputBandWidthChart.Detail = "兆字节"
	outputBandWidthChart.OnReload(func() {
		// 带宽
		{
			bandWidth := SharedLogger().OutputBandWidth()
			outputBandWidthChart.Value = float64(float64(bandWidth) / 1024 / 1024)
			outputBandWidthChart.Unit = "MB"

			max := math.Ceil(outputBandWidthChart.Value/float64(10)) * float64(10)
			if max == 0 {
				max = 20
			}
			outputBandWidthChart.Max = max
		}
	})
	widget.AddChart(outputBandWidthChart)

	// 日志
	chart := teacharts.NewTable()
	chart.Name = "即时日志"
	chart.OnReload(func() {
		// 日志
		accessLogs := SharedLogger().ReadNewLogs("", 5)
		chart.ResetRows()
		for _, accessLog := range accessLogs {
			chart.AddRow("<em>" + accessLog.TimeLocal + " @" + accessLog.Host + "</em><br/> <span>\"" + accessLog.Request + "\"</span>")
		}
	})
	widget.AddChart(chart)

	// 添加Widget
	plugin.AddWidget(widget)
}
