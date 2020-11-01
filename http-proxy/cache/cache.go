package cache

import (
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

type ProxyCache struct {
	data *cache.Cache
	path string
}

func InitCache(path string) ProxyCache {
	c := cache.New(30*time.Second, 30*time.Second)
	var proxyCache ProxyCache
	proxyCache.data = c
	proxyCache.path = path
	return proxyCache
}

func (c ProxyCache) Set(key string, raw []byte)  {
	if strings.HasPrefix(key, c.path) {
		c.data.SetDefault(key, raw)
	}
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