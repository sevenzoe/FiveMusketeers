package controllers

import (
	"github.com/dchest/captcha"
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisPool struct {
	pool redis.Pool
}

type captchaStore struct {
	store redis.Pool
}

var gCaptchaExpireTime = time.Duration(10) * time.Minute

func NewRedisPool(redisHost string, maxIdle, maxActive, redisDB int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}

			if redisDB != -1 {
				_, err = c.Do("SELECT", redisDB)
				if err != nil {
					return nil, err
				}
			}

			return c, err
		},
	}
}

// type Store interface
func NewCaptchaStore(redisHost string, maxIdle, maxActive, redisDB int) captcha.Store {
	cs := new(captchaStore)
	cs.store = redis.Pool{
		MaxIdle:   int(maxIdle),
		MaxActive: int(maxActive), // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}

			if redisDB != -1 {
				_, err = c.Do("SELECT", redisDB)
				if err != nil {
					return nil, err
				}
			}

			return c, err
		},
	}

	return cs
}

// 实现captcha.Store 接口方法
// Set sets the digits for the captcha id.
func (cs *captchaStore) Set(id string, digits []byte) {
	conn := cs.store.Get()
	defer conn.Close()

	_, err := conn.Do("SET", id, digits)
	if err != nil {

	}

	_, err = conn.Do("EXPIRE", gCaptchaExpireTime)
	if err != nil {

	}
}

// 实现captcha.Store 接口方法
// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (cs *captchaStore) Get(id string, clear bool) (digits []byte) {
	conn := cs.store.Get()
	defer conn.Close()

	value, err := conn.Do("GET", id)
	if err != nil {
		// TODO:
	}

	digits, err = redis.Bytes(value, err)
	if !clear {
		return
	}

	_, err = conn.Do("DEL", id)
	if err != nil {
		// TODO:
	}

	return
}
