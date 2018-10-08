package apps

import (
	"github.com/iwind/TeaWebCode/teaservices/probes"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teaplugins"
	"github.com/iwind/TeaWebCode/teacharts"
	"fmt"
)

type PHPFPMProbe struct {
	probes.Probe
}

func (this *PHPFPMProbe) Run() {
	this.InitOnce(func() {
		logs.Println("probe php-fpm")

		widget := teaplugins.NewWidget()
		widget.Dashboard = true
		widget.Group = teaplugins.WidgetGroupService
		widget.Name = "PHP-FPM"
		this.Plugin.AddWidget(widget)

		widget.OnForceReload(func() {
			this.Run()
		})
	})

	widget := this.Plugin.Widgets[0]
	result := ps("php-fpm", []string{"php-fpm(:|$)"}, true)
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
