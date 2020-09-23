package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisDB redis client
var RedisDB = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDRESS"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB: func() int {
		db, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 0)
		if err != nil {
			panic(err)
		}
		return int(db)
	}(),
})

// GetToken read k
func GetToken(k string) (string, error) {
	val, err := RedisDB.Get(ctx, k).Result()
	if err != nil {
		log.Printf("Redis get failed:%v", err.Error())
		return "", err
	}
	return val, err
}

// SetToken set k = v with x expiration
func SetToken(k string, v string, x time.Duration) bool {
	err := RedisDB.Set(ctx, k, v, x).Err()
	if err != nil {
		log.Printf("Redis set failed:%v", err.Error())
		panic(err)
	}
	return true
}

// SetTokenIfNotExists set k = v, if k not already in redis
// with x expiration
func SetTokenIfNotExists(k string, v string, x time.Duration) bool {
	val, err := RedisDB.SetNX(ctx, k, v, x).Result()
	if err != nil {
		log.Printf("Redis set failed:%v", err.Error())
		panic(err)
	}
	return val
}
