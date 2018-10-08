package tealogs

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaWebCode/teamongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/iwind/TeaGo/lists"
	"time"
	"sync"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/iwind/TeaGo/timers"
)

var (
	accessLogger = NewAccessLogger()
)

// 访问日志记录器
type AccessLogger struct {
	queue chan *AccessLogItem

	logs            []*AccessLogItem
	timestamp       int64
	qps             int
	outputBandWidth int64
	inputBandWidth  int64

	collectionCacheMap map[string]*mongo.Collection
	processors         []Processor
}

type AccessLogItem struct {
	log *AccessLog
}

func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue:              make(chan *AccessLogItem, 10240),
		collectionCacheMap: map[string]*mongo.Collection{},
	}

	go logger.wait()
	return logger
}

func SharedLogger() *AccessLogger {
	return accessLogger
}

func (this *AccessLogger) Push(log *AccessLog) {
	this.queue <- &AccessLogItem{
		log: log,
	}
}

func (this *AccessLogger) client() *mongo.Client {
	return teamongo.SharedClient()
}

func (this *AccessLogger) collection() *mongo.Collection {
	collName := "logs." + timeutil.Format("Ymd")
	coll, found := this.collectionCacheMap[collName]
	if found {
		return coll
	}

	// 构建索引
	coll = this.client().Database("teaweb").Collection(collName)
	indexes := coll.Indexes()
	indexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.NewDocument(bson.EC.Int32("status", 1), bson.EC.Int32("timestamp", 1)),
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})
	indexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.NewDocument(bson.EC.Int32("remoteAddr", 1), bson.EC.Int32("serverId", 1)),
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})

	this.collectionCacheMap[collName] = coll

	return coll
}

func (this *AccessLogger) wait() {
	timestamp := time.Now().Unix()

	var docs = []interface{}{}
	var docsLocker = sync.Mutex{}

	// 写入到数据库
	timers.Loop(1*time.Second, func(looper *timers.Looper) {
		// 写入到本地数据库
		if this.client() != nil {
			docsLocker.Lock()
			if len(docs) == 0 {
				docsLocker.Unlock()
				return
			}
			newDocs := docs
			docs = []interface{}{}
			docsLocker.Unlock()

			// 分析
			for _, doc := range newDocs {
				doc.(*AccessLog).Parse()
				doc.(*AccessLog).Id = objectid.New()

				// 其他处理器
				if len(this.processors) > 0 {
					for _, processor := range this.processors {
						processor.Process(doc.(*AccessLog))
					}
				}
			}

			total := len(newDocs)

			// 批量写入数据库
			bulkSize := 1024
			offset := 0
			for {
				end := offset + bulkSize
				if end > total {
					end = total
				}

				logs.Println("dump", end-offset, "docs ...")
				_, err := this.collection().InsertMany(context.Background(), newDocs[offset:end])
				if err != nil {
					logs.Error(err)
					return
				}
				logs.Println("done")

				offset = end
				if end >= total {
					break
				}
			}
		}
	})

	// 接收日志
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

		docsLocker.Lock()
		docs = append(docs, log)
		docsLocker.Unlock()
	}
}

// 添加处理器
func (this *AccessLogger) AddProcessor(processor Processor) {
	this.processors = append(this.processors, processor)
}

// 关闭
func (this *AccessLogger) Close() {
	if this.client() != nil {
		this.client().Disconnect(context.Background())
	}
}

// 读取日志
func (this *AccessLogger) ReadNewLogs(fromId string, size int64) []AccessLog {
	if this.client() == nil {
		return []AccessLog{}
	}

	if size <= 0 {
		size = 10
	}

	result := []AccessLog{}
	coll := this.collection()

	filter := map[string]interface{}{}
	if len(fromId) > 0 {
		objectId, err := objectid.FromHex(fromId)
		if err == nil {
			filter["_id"] = map[string]interface{}{
				"$gt": objectId,
			}
		} else {
			logs.Error(err)
		}
	}

	opts := []findopt.Find{}
	isReverse := false

	if len(fromId) == 0 {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("_id", -1))))
		opts = append(opts, findopt.Limit(size))
		isReverse = true
	} else {
		opts = append(opts, findopt.Sort(bson.NewDocument(bson.EC.Int32("_id", 1))))
		opts = append(opts, findopt.Limit(size))
	}

	cursor, err := coll.Find(context.Background(), filter, opts ...)
	if err != nil {
		logs.Error(err)
		return []AccessLog{}
	}
	defer cursor.Close(context.Background())

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
	coll := this.collection()
	filter := bson.NewDocument(
		bson.EC.SubDocument("status", bson.NewDocument(bson.EC.Int64("$lt", 400))),
		bson.EC.SubDocument("timestamp", bson.NewDocument(bson.EC.Int64("$lte", toTimestamp), bson.EC.Int64("$gte", fromTimestamp))),
	)
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		logs.Error(err)
		return 0
	}

	return count
}

func (this *AccessLogger) CountFailLogs(fromTimestamp int64, toTimestamp int64) int64 {
	coll := this.collection()
	filter := bson.NewDocument(
		bson.EC.SubDocument("status", bson.NewDocument(bson.EC.Int64("$gte", 400))),
		bson.EC.SubDocument("timestamp", bson.NewDocument(bson.EC.Int64("$lte", toTimestamp), bson.EC.Int64("$gte", fromTimestamp))),
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
