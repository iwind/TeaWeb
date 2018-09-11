package log

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaWebCode/teaweb/helpers"
	"time"
)

type GetAction actions.Action

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
	this.Data["qps"] = logger.QPS()

	this.Success()
}
