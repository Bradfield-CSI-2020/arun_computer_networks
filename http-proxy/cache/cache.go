package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type ProxyCache struct {
	data *cache.Cache
}

func InitCache() ProxyCache {
	c := cache.New(30*time.Second, 30*time.Second)
	var proxyCache ProxyCache
	proxyCache.data = c
	return proxyCache
}


func (c ProxyCache) Set(key string, raw []byte)  {
	c.data.SetDefault(key, raw)
}

func (c ProxyCache) Get(key string) []byte {
	value, _, found := c.data.GetWithExpiration(key)

	if found {
		// todo: what does this mean?
		return value.([]byte)
	} else {
		return nil
	}

}