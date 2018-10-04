package teastats

import (
	"github.com/iwind/TeaGo/utils/time"
	"github.com/iwind/TeaWebCode/tealogs"
	"strings"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"github.com/iwind/TeaGo/logs"
	"time"
	"github.com/iwind/TeaGo/types"
)

type HourlyUVStat struct {
	ServerId string `bson:"serverId" json:"serverId"` // 服务ID
	Hour     string `bson:"hour" json:"hour"`         // 小时，格式为：YmdH
	Count    int64  `bson:"count" json:"count"`       // 数量
}

func (this *HourlyUVStat) Init() {
	coll := findCollection("stats.uv.hourly", nil)
	coll.CreateIndex(map[string]bool{
		"hour": true,
	})
	coll.CreateIndex(map[string]bool{
		"hour":     true,
		"serverId": true,
	})
}

func (this *HourlyUVStat) Process(accessLog *tealogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	hour := timeutil.Format("YmdH")

	// 是否已存在
	result := findCollection("logs."+timeutil.Format("Ymd"), nil).FindOne(context.Background(), bson.NewDocument(bson.EC.String("remoteAddr", accessLog.RemoteAddr), bson.EC.String("serverId", accessLog.ServerId)), findopt.Projection(map[string]int{
		"id": 1,
	}))

	existAccessLog := map[string]interface{}{}
	if result.Decode(existAccessLog) != mongo.ErrNoDocuments {
		return
	}

	coll := findCollection("stats.uv.hourly", this.Init)

	stat := bson.NewDocument(
		bson.EC.SubDocument("$set", bson.NewDocument(
			bson.EC.String("serverId", accessLog.ServerId),
			bson.EC.String("hour", hour),
		)),
		bson.EC.SubDocument("$inc", bson.NewDocument(
			bson.EC.Int64("count", 1),
		)),
	)

	_, err := coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("hour", hour),
	), stat, updateopt.OptUpsert(true))
	if err != nil {
		logs.Error(err)
	}
}

func (this *HourlyUVStat) ListLatestHours(hours int) []map[string]interface{} {
	if hours <= 0 {
		hours = 24
	}

	result := []map[string]interface{}{}
	for i := hours - 1; i >= 0; i -- {
		hour := timeutil.Format("YmdH", time.Now().Add(time.Duration(-i)*time.Hour))
		total := this.SumHourUV([]string{hour})
		result = append(result, map[string]interface{}{
			"hour":  hour,
			"total": total,
		})
	}
	return result
}

func (this *HourlyUVStat) SumHourUV(hours []string) int64 {
	if len(hours) == 0 {
		return 0
	}
	sumColl := findCollection("stats.uv.hourly", nil)
	sumCursor, err := sumColl.Aggregate(context.Background(), bson.NewArray(bson.VC.DocumentFromElements(
		bson.EC.SubDocumentFromElements(
			"$match",
			bson.EC.Interface("hour", map[string]interface{}{
				"$in": hours,
			}),
		),
	), bson.VC.DocumentFromElements(bson.EC.SubDocumentFromElements(
		"$group",
		bson.EC.Interface("_id", nil),
		bson.EC.Interface("total", map[string]interface{}{
			"$sum": "$count",
		}),
	))))
	if err != nil {
		logs.Error(err)
		return 0
	}
	defer sumCursor.Close(context.Background())

	if sumCursor.Next(context.Background()) {
		sumMap := map[string]interface{}{}
		err = sumCursor.Decode(sumMap)
		if err == nil {
			return types.Int64(sumMap["total"])
		} else {
			logs.Error(err)
		}
	}

	return 0
}
