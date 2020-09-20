package main

import (
	"context"
)

var ctx = context.Background()

// SetToken haha
func SetToken(k string, v string) bool {
	err := RedisDB.Set(ctx, k, v, 0).Err()
	if err != nil {
		panic(err)
	}
	return true
}

// GetToken haha
func GetToken(k string) string {
	val, err := RedisDB.Get(ctx, k).Result()
	if err != nil {
		panic(err)
	}
	return val
}
