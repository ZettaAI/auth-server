package main

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

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
