// Author: huaxr
// Time:   2021/12/10 下午2:38
// Git:    huaxr

package kv

import (
	"context"
	"github.com/huaxr/magicflow/component/logx"

	"testing"
)

func TestC(t *testing.T) {
	ctx := context.Background()
	redisKey := "aaa"
	redisCli := RedisKv.GetCli()
	redisCli.RPush(ctx, redisKey, "1", "@", "3")
	//c := redisCli.Incr(ctx, redisKey).Val()
	//redisCli.Expire(ctx, redisKey, 5*time.Second)
	res := redisCli.LRange(ctx, redisKey, 0, -1)
	a, err := res.Result()
	logx.L().Infof("xx", a, err)

	redisKey2 := "bbb"
	redisCli.HSet(ctx, redisKey2, "1", "2", "3", "1")
	xx := redisCli.HGetAll(ctx, redisKey2)
	b, err := xx.Result()
	logx.L().Infof("xx", b, err)

	redisKey3 := "ccc"
	redisCli.SAdd(ctx, redisKey3, "1", "2", "3", "1")
	c, err := redisCli.SMembers(ctx, redisKey3).Result()
	logx.L().Infof("xx", c, err)

}
