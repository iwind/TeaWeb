package tealog

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaWebCode/teamongo"
	"context"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/iwind/TeaGo/lists"
)

var (
	accessLogger = NewAccessLogger()
)

type AccessLogger struct {
	queue  chan *AccessLogItem
	client *mongo.Client
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
	lastId := int64(0)

	for {
		item := <-this.queue
		log := item.log

		newId := time.Now().UnixNano()
		if newId == lastId {
			time.Sleep(1 * time.Nanosecond)
			newId = time.Now().UnixNano()
		}
		lastId = newId
		log.Id = newId

		// 分析日志
		log.parse()

		// 写入到本地数据库
		// @TODO 批量写入
		if this.client != nil {
			logs.Println("write access log to db")
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

func (this *AccessLogger) CountSuccessLogs() int64 {
	coll := this.client.Database("teaweb").Collection("accessLogs")
	filter := bson.NewDocument(bson.EC.SubDocument("status", bson.NewDocument(bson.EC.Int64("$lt", 400))))
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		logs.Error(err)
		return 0
	}

	return count
}

func (this *AccessLogger) CountFailLogs() int64 {
	coll := this.client.Database("teaweb").Collection("accessLogs")
	filter := bson.NewDocument(bson.EC.SubDocument("status", bson.NewDocument(bson.EC.Int64("$gte", 400))))
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		logs.Error(err)
		return 0
	}

	return count
}
