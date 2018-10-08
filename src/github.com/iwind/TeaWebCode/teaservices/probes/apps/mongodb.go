package apps

import (
	"github.com/iwind/TeaWebCode/teaservices/probes"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teaplugins"
	"github.com/iwind/TeaWebCode/teacharts"
	"fmt"
)

type MongoDBProbe struct {
	probes.Probe
}

func (this *MongoDBProbe) Run() {
	this.InitOnce(func() {
		logs.Println("probe mongodb")

		widget := teaplugins.NewWidget()
		widget.Dashboard = true
		widget.Group = teaplugins.WidgetGroupService
		widget.Name = "MongoDB"
		this.Plugin.AddWidget(widget)

		widget.OnForceReload(func() {
			this.Run()
		})
	})

	widget := this.Plugin.Widgets[0]
	result := ps("mongod", []string{"mongod$"}, true)
	widget.ResetCharts()
	if len(result) == 0 {
		return
	}
	for _, proc := range result {
		chart := teacharts.NewTable()
		chart.AddRow("PID:", fmt.Sprintf("%d", proc.Pid), "<i class=\"ui icon dot circle green\"></i>")
		chart.SetWidth(25, 60, 15)
		widget.AddChart(chart)
	}
}
