package redis

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis"
	"imsdk/pkg/app"
	"time"
)

type config struct {
	Addr         string `toml:"addr"`
	Password     string `toml:"password"`
	Db           int    `toml:"dao"`
	PoolSize     int    `toml:"pool_size"`
	MinIdleConns int    `toml:"min_idle_conns"`
	IsEnableTls  int    `toml:"is_enable_tls"`
}

var (
	// Client redis连接资源
	Client *redis.Client
	conf   config
	NilErr = redis.Nil
)

// Start 启动redis
func Start() {
	err := app.Config().Bind("db", "redis", &conf)
	fmt.Println("redis config:", err, conf)
	if err == app.ErrNodeNotExists {
		return
	}
	opt := &redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.Db,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
	}
	if conf.Password != "" {
		opt.Password = conf.Password
	}
	if conf.IsEnableTls == 1 {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
	}
	Client = redis.NewClient(opt)
}

// CacheGet 获取指定key的值,如果值不存在,就执行f方法将返回值存入redis
func CacheGet(key string, expiration time.Duration, f func() string) string {
	cmd := Client.Get(key)
	var val string
	result, _ := cmd.Result()
	if len(result) == 0 {
		Client.Set(key, f(), expiration)
		return val
	}
	return result
}

func Lock(key string, expire int) bool {
	lockName := "pay:lock:" + key
	lockTimeOut := time.Duration(expire) * time.Second
	if ok, err := Client.SetNX(lockName, 1, lockTimeOut).Result(); err != nil && err != redis.Nil {
		return false
	} else if ok {
		return true
	} else if Client.TTL(lockName).Val() == -1 { // -2:失效；-1：无过期；
		Client.Expire(lockName, lockTimeOut)
	}
	return false
}

func Unlock(key string) bool {
	lockName := "pay:lock:" + key
	num, err := Client.Del(lockName).Result()
	if err != nil && err == redis.Nil {
		return true
	} else if err == nil && num > 0 {
		return true
	}
	return false
}
