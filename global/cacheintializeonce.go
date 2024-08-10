package global

import (
	"sync"

	"github.com/codecrafters-io/redis-starter-go/mapstorage"
)

var (
	cacheInstance *mapstorage.Cache
	once sync.Once
)

func InitCacheDB() {
	once.Do(func() {
		cacheInstance = mapstorage.NewCache()
	})
}

func GetCache() *mapstorage.Cache {
	return cacheInstance
}

