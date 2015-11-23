package helper

import (
	"github.com/astaxie/beego/cache"
	"github.com/garyburd/redigo/redis"
)

var GlobalCache cache.Cache

// 缓存操作.GET
func GetCacheInt(key string) int {
	var err error
	v, _ := redis.Int(GlobalCache.Get(key), err)
	return v
}

// 缓存操作.SET
func SetCacheInt(key string, val int, timeout int64) error {
	return GlobalCache.Put(key, val, timeout)
}

// 缓存操作.DELETE
func DelCacheInt(key string) error {
	return GlobalCache.Delete(key)
}

// 缓存操作.INCREASE
func IncCacheInt(key string) error {
	return GlobalCache.Incr(key)
}

func GetCacheString(key string) string {
	var err error
	v, _ := redis.String(GlobalCache.Get(key), err)
	return v
}

func SetCacheString(key string, val string, timeout int64) error {
	return GlobalCache.Put(key, val, timeout)
}

func DelCacheString(key string) error {
	return GlobalCache.Delete(key)
}
