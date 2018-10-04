package teastats

import (
	"github.com/iwind/TeaWebCode/tealogs"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type TopRequestStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份
	URL      string `bson:"url" json:"url"`           // URL
	Count    int64  `bson:"count" json:"count"`       // 耗时
}

func (this *TopRequestStat) Init() {
	coll := findCollection("stats.top.requests.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"region":   true,
		"url":      true,
	})
	coll.CreateIndex(map[string]bool{
		"count": false,
	})
}

func (this *TopRequestStat) Process(accessLog *tealogs.AccessLog) {
	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.requests.monthly", this.Init)

	url := accessLog.Scheme + "://" + accessLog.Host + accessLog.RequestURI

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("url", url),
			bson.EC.String("month", month),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("url", url),
		bson.EC.String("month", month),
	), stat, updateopt.OptUpsert(true))
}
