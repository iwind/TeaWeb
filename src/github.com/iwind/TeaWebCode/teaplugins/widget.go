package teaplugins

import (
	"github.com/iwind/TeaWebCode/teacharts"
	"github.com/iwind/TeaGo/utils/string"
	"sync"
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

	// 刷新回调
	onReloadFuncs []func()
	reloadLocker  sync.Mutex

	onForceReloadFuncs []func()
	forceReloadLocker  sync.Mutex
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

func (this *Widget) ResetCharts() {
	this.Charts = []teacharts.ChartInterface{}
}

func (this *Widget) OnReload(f func()) {
	this.reloadLocker.Lock()
	defer this.reloadLocker.Unlock()

	this.onReloadFuncs = append(this.onReloadFuncs, f)
}

func (this *Widget) Reload() {
	if len(this.onReloadFuncs) == 0 && len(this.Charts) == 0 {
		return
	}

	// 异步执行
	if len(this.onReloadFuncs) == 0 {
		go func() {
			for _, chart := range this.Charts {
				chart.Reload()
			}
		}()
	} else {
		go func() {
			for _, f := range this.onReloadFuncs {
				f()
			}

			for _, chart := range this.Charts {
				chart.Reload()
			}
		}()
	}
}

func (this *Widget) OnForceReload(f func()) {
	this.forceReloadLocker.Lock()
	defer this.forceReloadLocker.Unlock()

	this.onForceReloadFuncs = append(this.onForceReloadFuncs, f)
}

func (this *Widget) ForceReload() {
	if len(this.onForceReloadFuncs) == 0 && len(this.Charts) == 0 {
		return
	}

	// 异步执行
	if len(this.onForceReloadFuncs) == 0 {
		go func() {
			for _, chart := range this.Charts {
				chart.Reload()
			}
		}()
	} else {
		go func() {
			for _, f := range this.onForceReloadFuncs {
				f()
			}

			for _, chart := range this.Charts {
				chart.Reload()
			}
		}()
	}
}
