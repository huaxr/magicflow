// Author: huaxr
// Time:   2021/9/24 下午6:57
// Git:    huaxr

package kv

import (
	"context"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/go-redis/redis/v8"
	"runtime"
	"time"
)

var RedisKv *redisKv

type redisKv struct {
	cli *redis.Client
}

func LunchRedis(ctx context.Context) {
	v := confutil.GetConf().Db.Redis
	client := redis.NewClient(&redis.Options{
		Addr:        v.Host,
		Password:    v.Password, // no password set
		DB:          v.Db,       // use default DB
		PoolSize:    v.Poolsize,
		IdleTimeout: time.Second * time.Duration(v.Idletimeout),
		ReadTimeout: time.Second * time.Duration(v.Readtimeout),
		MaxRetries:  v.MaxRetry,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		logx.L().Errorf("init redis err:%v", err)
		return
	}
	RedisKv = &redisKv{
		cli: client,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				logx.L().Warnf("redis close")
				RedisKv.cli.Close()
				runtime.Goexit()
			}
		}
	}()
}

func (r redisKv) GetCli() *redis.Client {
	return r.cli
}

func (r redisKv) KVGet(key string) (val interface{}, exist bool) {
	res := RedisKv.cli.Get(context.Background(), key)
	if res == nil {
		return nil, false
	}
	return res, true
}

func (r redisKv) KVSet(key string, val interface{}, duration time.Duration) {
	if RedisKv == nil {
		return
	}
	RedisKv.cli.Set(context.Background(), key, val, duration)
}

func (r redisKv) use() {
	//redis.SAdd(ctx, key, n.NodeCode)
	//redis.Expire(ctx, key, duration)
	//mem, err := redis.SMembers(ctx, key).Result()
}
