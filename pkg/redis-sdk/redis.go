package redissdk

import (
	"context"
	"log"
	"time"

	"github.com/Zeroaril7/nobi-technical-test/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {

	var err error
	for retries := 0; retries < 5; retries++ {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     config.Config("REDIS_URL"),
			Password: config.Config("REDIS_PASSWORD"),
			DB:       0,
		})
		_, err = RedisClient.Ping(Ctx).Result()
		if err == nil {
			break
		}
		log.Println("Retrying Redis connection...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
