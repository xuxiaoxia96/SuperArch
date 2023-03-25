package redis

import (
	"SuperArch/conf"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)


func GetRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Cfg.Redis.Host, conf.Cfg.Redis.Port),
		Password: conf.Cfg.Redis.Password,
		DB:       int(conf.Cfg.Redis.DB),
		PoolSize: 100, // Connection Pool
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}

func InsertHset(hashKey string, kv map[string]interface{}){
	ctx := context.Background()
	rdb, err := GetRedisClient()
	if err != nil{
		logrus.Errorf("[Redis Module][InsertHset][InitRedisClient] %s", err)
		return
	}

	rdb.HSet(ctx, hashKey, kv)
}

func UpdateHset(hashKey string, kv map[string]interface{}){
	ctx := context.Background()
	rdb, err := GetRedisClient()
	if err != nil{
		logrus.Errorf("[Redis Module][UpdateHset][InitRedisClient] %s", err)
		return
	}

	rdb.HSet(ctx, hashKey, kv)
}

func DeleteHsetByKey(hashKey string, k string){
	ctx := context.Background()
	rdb, err := GetRedisClient()
	if err != nil{
		logrus.Errorf("[Redis Module][DeleteHset][InitRedisClient] %s", err)
		return
	}

	rdb.HDel(ctx, hashKey, k)
}

func GetALLFromHset(hashKey string) map[string]string{
	ctx := context.Background()
	rdb, err := GetRedisClient()
	if err != nil{
		logrus.Errorf("[Redis Module][GetAllFromHset][InitRedisClient] %s", err)
		return nil
	}

	kv := rdb.HGetAll(ctx, hashKey)

	return kv.Val()
}

func GetValFromHsetBykey(hashKey, k string) string{
	ctx := context.Background()
	rdb, err := GetRedisClient()
	if err != nil{
		logrus.Errorf("[Redis Module][GetValFromHsetBykey][InitRedisClient] %s", err)
		return ""
	}

	res,err := rdb.HGet(ctx, hashKey, k).Result()
	if err != nil{
		logrus.Errorf("[Redis Module][GetValFromHsetBykey][Get Hset Val] %s", err)
		return ""
	}

	return res
}
