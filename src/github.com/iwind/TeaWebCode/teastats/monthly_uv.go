package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/iwind/TeaWebCode/tealog"
	"strings"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

type MonthlyUVStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Month    string `bson:"month" json:"month"`       // 月份，格式为：Ym
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *MonthlyUVStat) Init() {
	coll := findCollection("stats.uv.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
	coll.CreateIndex(map[string]bool{
		"month":    true,
		"serverId": true,
	})
}

func (this *MonthlyUVStat) Process(accessLog *tealog.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	month := timeutil.Format("Ym")

	// 是否已存在
	result := findCollection("log."+timeutil.Format("Ymd"), nil).FindOne(context.Background(), bson.NewDocument(bson.EC.String("remoteAddr", accessLog.RemoteAddr), bson.EC.String("serverId", accessLog.ServerId)), findopt.Projection(map[string]int{
		"id": 1,
	}))

	existAccessLog := map[string]interface{}{}
	if result.Decode(existAccessLog) != mongo.ErrNoDocuments {
		return
	}

	coll := findCollection("stats.uv.monthly", this.Init)

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
