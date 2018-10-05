package teaplugins

import (
	"github.com/iwind/TeaWebCode/teacharts"
	"github.com/iwind/TeaGo/utils/string"
)

type Widget struct {
	Name  string `json:"name"`  // 名称
	Icon  string `json:"icon"`  // Icon @TODO
	Title string `json:"title"` // 标题

	URL       string `json:"url"`       // 外部链接URL
	MoreURL   string `json:"moreUrl"`   // 更多信息链接URL
	TopBar    bool   `json:"topBar"`    // 是否顶部工具栏可用
	MenuBar   bool   `json:"menuBar"`   // 是否菜单栏可用
	HelperBar bool   `json:"helperBar"` // 是否小助手栏可用
	Dashboard bool   `json:"dashboard"` // 是否在仪表盘可用

	Group WidgetGroup `json:"group"`

	// 图表类型
	Charts []teacharts.ChartInterface `json:"charts"`
}

func NewWidget() *Widget {
	return &Widget{
		Charts: []teacharts.ChartInterface{},
	}
}

func (this *Widget) AddChart(chart teacharts.ChartInterface) {
	if len(chart.UniqueId()) == 0 {
		chart.SetUniqueId(stringutil.Rand(16))
	}
	this.Charts = append(this.Charts, chart)
}
