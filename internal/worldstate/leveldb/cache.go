package leveldb

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger-labs/orion-server/pkg/types"
)

var keySep = []byte{0x00}

// cache holds both the system and user cache
type cache struct {
	dataCache *fastcache.Cache
}

// newCache creates a Cache. The size of the data cache, in terms of MB, is
// specified via cacheSize parameter. Note that the maximum memory consumption of fastcache
// would be in the multiples of 32 MB (due to 512 buckets & an equal number of 64 KB chunks per bucket).
// If the cacheSizeMBs is not a multiple of 32 MB, the fastcache would round the size
// to the next multiple of 32 MB.
func newCache(dataCacheSizeMBs int) *cache {
	cache := &cache{}

	cache.dataCache = fastcache.New(dataCacheSizeMBs * 1024 * 1024)
	return cache
}

// getState returns the value for a given namespace and key from
// the cache
func (c *cache) getState(namespace, key string) (*types.ValueWithMetadata, error) {
	cacheKey := constructCacheKey(namespace, key)

	if !c.dataCache.Has(cacheKey) {
		return nil, nil
	}
	cacheValue := &types.ValueWithMetadata{}
	valBytes := c.dataCache.Get(nil, cacheKey)
	if err := proto.Unmarshal(valBytes, cacheValue); err != nil {
		return nil, err
	}
	return cacheValue, nil
}

// putState stores a given value in the cache
func (c *cache) putState(namespace, key string, cacheValue []byte) error {
	cacheKey := constructCacheKey(namespace, key)

	if c.dataCache.Has(cacheKey) {
		c.dataCache.Del(cacheKey)
	}

	c.dataCache.Set(cacheKey, cacheValue)
	return nil
}

// putStateIfExist stores a given value in the cache only if the key already
// exists in the cache
func (c *cache) putStateIfExist(namespace, key string, cacheValue []byte) {
	cacheKey := constructCacheKey(namespace, key)

	if c.dataCache.Has(cacheKey) {
		c.dataCache.Del(cacheKey)
		c.dataCache.Set(cacheKey, cacheValue)
	}
}

func (c *cache) delState(namespace, key string) {
	cacheKey := constructCacheKey(namespace, key)

	if c.dataCache.Has(cacheKey) {
		c.dataCache.Del(cacheKey)
	}
}

// Reset removes all the items from the cache.
func (c *cache) Reset() {
	c.dataCache.Reset()
}

func constructCacheKey(namespace, key string) []byte {
	var cacheKey []byte
	cacheKey = append(cacheKey, []byte(namespace)...)
	cacheKey = append(cacheKey, keySep...)
	return append(cacheKey, []byte(key)...)
}
