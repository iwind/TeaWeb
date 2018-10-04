package teastats

import (
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type TopStateStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份
	Region   string `bson:"region" json:"region"`     // 地区
	State    string `bson:"state" json:"state"`       // 省份|州
	Count    int64  `bson:"count" json:"count"`       // 访问数量
}

func (this *TopStateStat) Init() {
	coll := findCollection("stats.top.states.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"state":    true,
		"region":   true,
		"month":    true,
	})
	coll.CreateIndex(map[string]bool{
		"count": false,
	})
}

func (this *TopStateStat) Process(accessLog *tealog.AccessLog) {
	if len(accessLog.Extend.Geo.Region) == 0 || len(accessLog.Extend.Geo.State) == 0 {
		return
	}
	region := accessLog.Extend.Geo.Region
	state := accessLog.Extend.Geo.State

	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.states.monthly", this.Init)

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("region", region),
			bson.EC.String("state", state),
			bson.EC.String("month", month),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("region", region),
		bson.EC.String("state", state),
		bson.EC.String("month", month),
	), stat, updateopt.OptUpsert(true))
}
