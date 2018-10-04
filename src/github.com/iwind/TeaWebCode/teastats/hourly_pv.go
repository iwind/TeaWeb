package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"strings"
)

type HourlyPVStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Hour     string `bson:"hour" json:"hour"`         // 小时，格式为：YmdH
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *HourlyPVStat) Init() {
	coll := findCollection("stats.pv.hourly", nil)
	coll.CreateIndex(map[string]bool{
		"hour": true,
	})
	coll.CreateIndex(map[string]bool{
		"hour":     true,
		"serverId": true,
	})
}

func (this *HourlyPVStat) Process(accessLog *tealog.AccessLog) {
	if !strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		return
	}

	hour := timeutil.Format("YmdH")
	coll := findCollection("stats.pv.hourly", this.Init)

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("hour", hour),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("hour", hour),
	), stat, updateopt.OptUpsert(true))
}
