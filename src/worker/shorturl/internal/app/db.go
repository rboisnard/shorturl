package app

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type ShortDB struct {
	DBClient *redis.Client
}

func (shortDB *ShortDB) Init() {
	shortDB.DBClient = redis.NewClient(&redis.Options{
		Addr:     Config.RedisIP + ":" + Config.RedisPort,
		Password: "",
		DB:       0,
	})
}

func (shortDB ShortDB) Ping() (string, error) {
	return shortDB.DBClient.Ping(context.TODO()).Result()
}
