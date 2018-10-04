package teastats

import (
	"github.com/iwind/TeaWebCode/teamongo"
	"sync"
	"github.com/iwind/TeaWebCode/tealog"
)

var collectionsMap = map[string]*teamongo.Collection{} // name => collection
var collectionsMutex = &sync.Mutex{}
var processors = []tealog.Processor{
	new(DailyPVStat),
	new(HourlyPVStat),
	new(MonthlyPVStat),

	new(DailyRequestsStat),
	new(HourlyRequestsStat),
	new(MonthlyRequestsStat),

	new(DailyUVStat),
	new(HourlyUVStat),
	new(MonthlyUVStat),

	new(TopRegionStat),
	new(TopStateStat),
	new(TopOSStat),
	new(TopBrowserStat),
	new(TopRequestStat),
	new(TopCostStat),
}

type Processor struct {
}

func (this *Processor) Process(accessLog *tealog.AccessLog) {
	for _, processor := range processors {
		processor.Process(accessLog)
	}
}

func findCollection(collectionName string, initFunc func()) *teamongo.Collection {
	collectionsMutex.Lock()
	defer collectionsMutex.Unlock()

	coll, found := collectionsMap[collectionName]
	if found {
		return coll
	}

	coll = teamongo.FindCollection(collectionName)
	collectionsMap[collectionName] = coll

	// 初始化
	if initFunc != nil {
		go initFunc()
	}

	return coll
}
