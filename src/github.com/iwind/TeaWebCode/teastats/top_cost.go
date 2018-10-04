package teastats

import (
	"github.com/iwind/TeaWebCode/tealogs"
	"github.com/iwind/TeaGo/utils/time"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/logs"
)

type TopCostStat struct {
	ServerId  string  `bson:"serverId" json:"serverId"`   // 服务ID
	Month     string  `bson:"month" json:"month"`         // 月份
	URL       string  `bson:"url" json:"url"`             // URL
	Cost      float64 `bson:"cost" json:"cost"`           // 平均耗时
	TotalCost float64 `bson:"totalCost" json:"totalCost"` // 总耗时
	Count     int64   `bson:"count" json:"count"`         // 请求数量
}

func (this *TopCostStat) Init() {
	coll := findCollection("stats.top.cost.monthly", nil)
	coll.CreateIndex(map[string]bool{
		"serverId": true,
		"region":   true,
		"url":      true,
	})
	coll.CreateIndex(map[string]bool{
		"cost": false,
	})
}

func (this *TopCostStat) Process(accessLog *tealogs.AccessLog) {
	month := timeutil.Format("Ym")
	coll := findCollection("stats.top.cost.monthly", this.Init)

	url := accessLog.Scheme + "://" + accessLog.Host + accessLog.RequestURI

	filter := map[string]interface{}{
		"serverId": accessLog.ServerId,
		"url":      url,
		"month":    month,
	}

	stat := map[string]interface{}{
		"$set": map[string]interface{}{
			"serverId": accessLog.ServerId,
			"url":      url,
			"month":    month,
			"cost":     accessLog.RequestTime,
		},
		"$inc": map[string]interface{}{
			"count":     1,
			"totalCost": accessLog.RequestTime,
		},
	}

	result := coll.FindOneAndUpdate(context.Background(), filter, stat, findopt.OptUpsert(true), findopt.Projection(map[string]int{
		"_id":       1,
		"totalCost": 1,
		"count":     1,
	}))

	m := map[string]interface{}{}
	if result.Decode(m) != mongo.ErrNoDocuments {
		count := types.Int64(m["count"]) + 1
		totalCost := types.Float64(m["totalCost"]) + accessLog.RequestTime
		avgCost := totalCost / float64(count)
		_, err := coll.UpdateOne(context.Background(), map[string]interface{}{
			"_id": m["_id"],
		}, map[string]interface{}{
			"$set": map[string]interface{}{
				"cost": avgCost,
			},
		})
		if err != nil {
			logs.Error(err)
		}
	}
}
