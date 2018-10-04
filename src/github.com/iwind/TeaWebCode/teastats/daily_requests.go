package teastats

import (
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaGo/utils/time"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type DailyRequestsStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Day      string `bson:"day" json:"day"`           // 日期，格式为：Ymd
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *DailyRequestsStat) Init() {
	coll := findCollection("stats.requests.daily", nil)
	coll.CreateIndex(map[string]bool{
		"day": true,
	})
	coll.CreateIndex(map[string]bool{
		"day":      true,
		"serverId": true,
	})
}

func (this *DailyRequestsStat) Process(accessLog *tealog.AccessLog) {
	day := timeutil.Format("Ymd")
	coll := findCollection("stats.requests.daily", this.Init)

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("day", day),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("day", day),
	), stat, updateopt.OptUpsert(true))
}
