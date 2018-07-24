package tealog

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaWebCode/teadb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaWebCode/teautils"
)

var (
	accessLogger = NewAccessLogger()
)

type AccessLogger struct {
	queue chan *AccessLogItem
	db    *teadb.DB
}

type AccessLogItem struct {
	log     *AccessLog
	writers []AccessLogWriter
}

func NewAccessLogger() *AccessLogger {
	logger := &AccessLogger{
		queue: make(chan *AccessLogItem, 10240),
	}

	accessDB, err := teadb.NewDB(Tea.LogFile("accesslog"))
	if err != nil {
		logs.Error(err)
	} else {
		logger.db = accessDB
	}

	go logger.wait()
	return logger
}

func SharedLogger() *AccessLogger {
	return accessLogger
}

func (logger *AccessLogger) Push(log *AccessLog, writers []AccessLogWriter) {
	logger.queue <- &AccessLogItem{
		log:     log,
		writers: writers,
	}
}

func (logger *AccessLogger) wait() {
	for {
		item := <-logger.queue
		log := item.log

		// 分析日志
		log.parse()

		// 写入到本地数据库
		if logger.db != nil {
			logger.db.PutStruct("ACCESSLOG", log)
		}

		// 输出日志到各个接口
		if len(item.writers) > 0 {
			for _, writer := range item.writers {
				writer.Write(log)
			}
		}
	}
}

func (logger *AccessLogger) Close() {
	if logger.db != nil {
		logger.db.Close()
	}
}

func (logger *AccessLogger) ReadLogs(fromId int64, size int64) []AccessLog {
	if logger.db == nil {
		return []AccessLog{}
	}

	if size <= 0 {
		size = 10
	}

	result := []AccessLog{}
	records, err := logger.db.NewQuery("ACCESSLOG").
		Reverse(true).
		Limit(size).
		FindAll()
	if err != nil {
		return result
	}

	for _, record := range records {
		accessLog := AccessLog{}
		err := teautils.MapToObjectJSON(record.Value, &accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}

		accessLog.Id = record.Id()
		result = append(result, accessLog)
	}

	return result
}
