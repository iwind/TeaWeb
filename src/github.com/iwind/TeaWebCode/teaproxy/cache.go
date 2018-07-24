package teaproxy

import (
	"runtime"
	"github.com/pbnjay/memory"
	"sync"
	"time"
)

type FixedCache struct {
	maxMemory uint64
	items     map[string]interface{}
	mutex     *sync.Mutex
}

func NewFixedCache() *FixedCache {
	cache := &FixedCache{
		maxMemory: memory.TotalMemory() / 8,
		items:     map[string]interface{}{},
		mutex:     &sync.Mutex{},
	}

	go cache.checkMemory()

	return cache
}

func (cache *FixedCache) Add(key string, object interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.items[key] = object
}

func (cache *FixedCache) Get(key string) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	object, found := cache.items[key]
	return object, found
}

func (cache *FixedCache) Trim() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	max := len(cache.items) / 2
	count := 0
	for key := range cache.items {
		if count < max {
			delete(cache.items, key)
		} else {
			break
		}

		count ++
	}
}

func (cache *FixedCache) checkMemory() {
	for {
		func() {
			stat := &runtime.MemStats{}
			runtime.ReadMemStats(stat)

			total := stat.TotalAlloc
			if total > cache.maxMemory {
				cache.Trim()
			}
		}()
		time.Sleep(5 * time.Second)
	}
}
