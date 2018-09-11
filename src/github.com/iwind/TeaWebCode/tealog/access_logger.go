package tealog

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaWebCode/teamongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"time"
)

var (
	accessLogger = NewAccessLogger()
)

// 访问日志记录器
type AccessLogger struct {
	queue  chan *AccessLogItem
	client *mongo.Client

	timestamp       int64
	qps             int
	outputBandWidth int64
	inputBandWidth  int64
}

type AccessLogItem struct {
	log     *AccessLog
	writers []AccessLogWriter
}

func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue:  make(chan *AccessLogItem, 10240),
		client: teamongo.SharedClient(),
	}

	go logger.wait()
	return logger
}

func SharedLogger() *AccessLogger {
	return accessLogger
}

func (this *AccessLogger) Push(log *AccessLog, writers []AccessLogWriter) {
	this.queue <- &AccessLogItem{
		log:     log,
		writers: writers,
	}
}

func (this *AccessLogger) wait() {
	latestDoc := this.client.Database("teaweb").
		Collection("accessLogs").
		FindOne(context.Background(), nil, findopt.Sort(bson.NewDocument(bson.EC.Int32("id", -1))))
	one := maps.Map{}
	err := latestDoc.Decode(one)

	var newId int64
	if err != nil {
		newId = time.Now().UnixNano() / 1000000
	} else {
		newId = one.GetInt64("id")
	}
	logs.Println("start log id:", newId)

	timestamp := time.Now().Unix()

	for {
		item := <-this.queue
		log := item.log

		// 计算QPS和BandWidth
		this.timestamp = log.Timestamp
		if log.Timestamp == timestamp {
			this.qps ++
			this.inputBandWidth += log.RequestLength
			this.outputBandWidth += log.BytesSent
		} else {
			this.qps = 1
			this.inputBandWidth = log.RequestLength
			this.outputBandWidth = log.BytesSent
			timestamp = log.Timestamp
		}

		newId ++
		log.Id = newId

		// 分析日志
		log.parse()

		// 写入到本地数据库
		// @TODO 批量写入
		if this.client != nil {
			this.client.
				Database("teaweb").
				Collection("accessLogs").
				InsertOne(context.Background(), log)
		}

		// 输出日志到各个接口
		if len(item.writers) > 0 {
			for _, writer := range item.writers {
				writer.Write(log)
			}
		}
	}
}

func (this *AccessLogger) Close() {
	if this.client != nil {
		this.client.Disconnect(context.Background())
	}
}

// 读取日志
func (this *AccessLogger) ReadNewLogs(fromId int64, size int64) []AccessLog {
	if this.client == nil {
		return []AccessLog{}
	}

	if size <= 0 {
		size = 10
	}

	result := []AccessLog{}
	coll := this.client.Database("teaweb").Collection("accessLogs")

	filter := bson.NewDocument(bson.EC.SubDocument("id", bson.NewDocument(bson.EC.Int64("$gt", fromId))))

	opts := []findopt.Find{}
	isReverse := false
	if fromId <= 0 {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("id", -1))))
		opts = append(opts, findopt.Limit(100))
		isReverse = true
	} else {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("id", 1))))
		opts = append(opts, findopt.Limit(size))
	}

	cursor, err := coll.Find(context.Background(), filter, opts ...)
	if err != nil {
		logs.Error(err)
		return []AccessLog{}
	}

	for cursor.Next(context.Background()) {
		accessLog := AccessLog{}
		err := cursor.Decode(&accessLog)
		if err != nil {
			logs.Error(err)
			return []AccessLog{}
		}
		result = append(result, accessLog)
	}

	if !isReverse {
		lists.Reverse(result)
	}
	return result
}

func (this *AccessLogger) CountSuccessLogs(fromTimestamp int64, toTimestamp int64) int64 {
	coll := this.client.Database("teaweb").Collection("accessLogs")
	filter := bson.NewDocument(
		bson.EC.SubDocument("status", bson.NewDocument(bson.EC.Int64("$lt", 400))),
		bson.EC.SubDocument("msec", bson.NewDocument(bson.EC.Int64("$lte", toTimestamp), bson.EC.Int64("$gte", fromTimestamp))),
	)
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		logs.Error(err)
		return 0
	}

	return count
}

func (this *AccessLogger) CountFailLogs(fromTimestamp int64, toTimestamp int64) int64 {
	coll := this.client.Database("teaweb").Collection("accessLogs")
	filter := bson.NewDocument(
		bson.EC.SubDocument("status", bson.NewDocument(bson.EC.Int64("$gte", 400))),
		bson.EC.SubDocument("msec", bson.NewDocument(bson.EC.Int64("$lte", toTimestamp), bson.EC.Int64("$gte", fromTimestamp))),
	)
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		logs.Error(err)
		return 0
	}

	return count
}

func (this *AccessLogger) QPS() int {
	if time.Now().Unix()-this.timestamp < 2 {
		return this.qps
	}
	return 0
}

func (this *AccessLogger) InputBandWidth() int64 {
	if time.Now().Unix()-this.timestamp < 2 {
		return this.inputBandWidth
	}
	return 0
}

func (this *AccessLogger) OutputBandWidth() int64 {
	if time.Now().Unix()-this.timestamp < 2 {
		return this.outputBandWidth
	}
	return 0
}
