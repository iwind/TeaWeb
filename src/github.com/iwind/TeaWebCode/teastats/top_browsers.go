package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/iwind/TeaWebCode/tealog"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type TopBrowserStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份
	Family   string `bson:"family" json:"family"`     // 浏览器
	Version  string `bson:"version" json:"version"`   // 版本
	Count    int64  `bson:"count" json:"count"`       // 访问数量
}

func (this *TopBrowserStat) Init() {
	coll := findCollection("stats.top.browsers.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"family":   true,
		"version":  true,
		"month":    true,
	})
	coll.CreateIndex(map[string]bool{
		"count": false,
	})
}

func (this *TopBrowserStat) Process(accessLog *tealog.AccessLog) {
	if len(accessLog.Extend.Client.Browser.Family) == 0 {
		return
	}
	family := accessLog.Extend.Client.Browser.Family
	version := accessLog.Extend.Client.Browser.Major

	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.browsers.monthly", this.Init)

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("family", family),
			bson.EC.String("version", version),
			bson.EC.String("month", month),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("family", family),
		bson.EC.String("version", version),
		bson.EC.String("month", month),
	), stat, updateopt.OptUpsert(true))
}
