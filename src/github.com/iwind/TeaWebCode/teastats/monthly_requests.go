package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type MonthlyRequestsStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份，格式为：Ym
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *MonthlyRequestsStat) Init() {
	coll := findCollection("stats.requests.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
	coll.CreateIndex(map[string]bool{
		"month":    true,
		"serverId": true,
	})
}

func (this *MonthlyRequestsStat) Process(accessLog *tealog.AccessLog) {
	month := timeutil.Format("Ym")
	coll := findCollection("stats.requests.monthly", this.Init)

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("month", month),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("month", month),
	), stat, updateopt.OptUpsert(true))
}
