package teastats

import (
	"github.com/iwind/TeaWebCode/tealogs"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/iwind/TeaGo/logs"
	"time"
)

type TopStateStat struct {
	ServerId string  `bson:"serverId" json:"serverId"` // 服务ID
	Month    string  `bson:"month" json:"month"`       // 月份
	Region   string  `bson:"region" json:"region"`     // 地区
	State    string  `bson:"state" json:"state"`       // 省份|州
	Count    int64   `bson:"count" json:"count"`       // 访问数量
	Percent  float64 `bson:"percent" json:"percent"`   // 比例
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
	coll.CreateIndex(map[string]bool{
		"month": true,
	})
}

func (this *TopStateStat) Process(accessLog *tealogs.AccessLog) {
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

	_, err := coll.UpdateOne(context.Background(), bson.NewDocument(
		bson.EC.String("serverId", accessLog.ServerId),
		bson.EC.String("region", region),
		bson.EC.String("state", state),
		bson.EC.String("month", month),
	), stat, updateopt.OptUpsert(true))
	if err != nil {
		logs.Error(err)
	}
}

func (this *TopStateStat) List(size int64) (result []TopStateStat) {
	if size <= 0 {
		size = 10
	}

	result = []TopStateStat{}

	// 最近两个月
	months := []string{}
	month1 := timeutil.Format("Ym")
	month2 := timeutil.Format("Ym", time.Now().AddDate(0, -1, 0))
	if month1 != month2 {
		months = append(months, month1, month2)
	} else {
		months = append(months, month1)
	}

	// 总请求数量
	totalRequests := new(MonthlyRequestsStat).SumMonthRequests(months)

	// 开始查找
	coll := findCollection("stats.top.states.monthly", nil)
	cursor, err := coll.Find(context.Background(), map[string]interface{}{
		"month": map[string]interface{}{
			"$in": months,
		},
	}, findopt.Sort(map[string]interface{}{
		"count": -1,
	}), findopt.Limit(size))
	if err != nil {
		logs.Error(err)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		one := TopStateStat{}
		err := cursor.Decode(&one)
		if err == nil {
			if totalRequests > 0 {
				one.Percent = float64(one.Count) / float64(totalRequests)
			} else {
				one.Percent = 0
			}

			result = append(result, one)
		} else {
			logs.Error(err)
		}
	}

	return
}
