package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"time"
	"math"
)

type GetAction actions.Action

var lastTotal = int64(0)
var lastTime = time.Now()

func (this *GetAction) Run(params struct {
	FromId int64 `alias:"fromId" default:"-1"`
	Size   int64 `default:"10"`
	Auth   *helpers.UserMustAuth
}) {
	logger := tealog.SharedLogger()
	accessLogs := logger.ReadNewLogs(params.FromId, params.Size)
	this.Data["logs"] = accessLogs

	fromTime := time.Now().Add(-24 * time.Hour)
	toTime := time.Now()

	countSuccess := logger.CountSuccessLogs(fromTime.Unix(), toTime.Unix())
	countFail := logger.CountFailLogs(fromTime.Unix(), toTime.Unix())
	total := countSuccess + countFail
	this.Data["countSuccess"] = countSuccess
	this.Data["countFail"] = countFail
	this.Data["total"] = total
	this.Data["qps"] = 0

	if lastTotal > 0 {
		countRequests := total - lastTotal
		if countRequests > 0 {
			this.Data["qps"] = int64(math.Ceil(float64(countRequests) / time.Since(lastTime).Seconds()))
		}
	}

	lastTotal = total
	lastTime = time.Now()

	this.Success()
}
