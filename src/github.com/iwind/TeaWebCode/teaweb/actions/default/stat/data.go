package stat

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/teastats"
	"github.com/iwind/TeaGo/types"
)

type DataAction actions.Action

func (this *DataAction) Run(params struct {
	Type  string `default:"pv"`    // 数据类型：uv|pv\req
	Range string `default:"daily"` // 时间范围，hourly|daily|monthly
}) {

	title := ""
	labels := []string{}
	data := []int64{}

	if params.Type == "uv" {
		if params.Range == "hourly" {
			title = "24小时UV统计"
			for _, stat := range new(teastats.HourlyUVStat).ListLatestHours(24) {
				labels = append(labels, types.String(stat["hour"])[8:])
				data = append(data, types.Int64(stat["total"]))
			}
		} else if params.Range == "daily" {
			title = "14日UV统计"
			for _, stat := range new(teastats.DailyUVStat).ListLatestDays(14) {
				labels = append(labels, types.String(stat["day"])[4:])
				data = append(data, types.Int64(stat["total"]))
			}
		} else if params.Range == "monthly" {
			title = "月UV统计"
			for _, stat := range new(teastats.MonthlyUVStat).ListLatestMonths(12) {
				labels = append(labels, types.String(stat["month"]))
				data = append(data, types.Int64(stat["total"]))
			}
		}
	} else if params.Type == "pv" {
		if params.Range == "hourly" {
			title = "24小时PV统计"
			for _, stat := range new(teastats.HourlyPVStat).ListLatestHours(24) {
				labels = append(labels, types.String(stat["hour"])[8:])
				data = append(data, types.Int64(stat["total"]))
			}
		} else if params.Range == "daily" {
			title = "14日PV统计"
			for _, stat := range new(teastats.DailyPVStat).ListLatestDays(14) {
				labels = append(labels, types.String(stat["day"])[4:])
				data = append(data, types.Int64(stat["total"]))
			}
		} else if params.Range == "monthly" {
			title = "月PV统计"
			for _, stat := range new(teastats.MonthlyPVStat).ListLatestMonths(12) {
				labels = append(labels, types.String(stat["month"]))
				data = append(data, types.Int64(stat["total"]))
			}
		}
	} else if params.Type == "req" {
		if params.Range == "hourly" {
			title = "24小时访问量统计"
			for _, stat := range new(teastats.HourlyRequestsStat).ListLatestHours(24) {
				labels = append(labels, types.String(stat["hour"])[8:])
				data = append(data, types.Int64(stat["total"]))
			}
		} else if params.Range == "daily" {
			title = "14日访问量统计"
			for _, stat := range new(teastats.DailyRequestsStat).ListLatestDays(14) {
				labels = append(labels, types.String(stat["day"])[4:])
				data = append(data, types.Int64(stat["total"]))
			}
		} else if params.Range == "monthly" {
			title = "月访问量统计"
			for _, stat := range new(teastats.MonthlyRequestsStat).ListLatestMonths(12) {
				labels = append(labels, types.String(stat["month"]))
				data = append(data, types.Int64(stat["total"]))
			}
		}
	}

	this.Data["title"] = title
	this.Data["labels"] = labels
	this.Data["data"] = data

	this.Success()
}
