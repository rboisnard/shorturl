package app

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type ShortDB struct {
	DBClient  *redis.Client
	DBContext context.Context
}

// global var
var DBInstance ShortDB

func (shortDB *ShortDB) Init() {
	shortDB.DBClient = redis.NewClient(&redis.Options{
		Addr:     Config.RedisIP + ":" + Config.RedisPort,
		Password: "",
		DB:       0,
	})
	shortDB.DBContext = context.Background()
}

func (shortDB ShortDB) Ping() (string, error) {
	return shortDB.DBClient.Ping(shortDB.DBContext).Result()
}

func (shortDB ShortDB) GetLongURL(surl string) (string, error) {
	var now string = strconv.FormatInt(time.Now().Unix(), 10)
	pipe := shortDB.DBClient.TxPipeline()
	lurl := pipe.HGet(shortDB.DBContext, surl, "longurl")
	pipe.HIncrBy(shortDB.DBContext, surl, "count", 1)
	pipe.HSet(shortDB.DBContext, surl, "last", now)
	_, err := pipe.Exec(shortDB.DBContext)
	return lurl.Val(), err
}

func (shortDB ShortDB) AddLongURL(surl, lurl string) (int64, error) {
	var now string = strconv.FormatInt(time.Now().Unix(), 10)
	return shortDB.DBClient.HSet(shortDB.DBContext, surl, []string{
		"longurl", lurl,
		"creation", now,
		"last", now,
		"count", "0",
	}).Result()
}

func (shortDB ShortDB) CheckLongURL(surl, lurl string) (bool, error) {
	rec, err := shortDB.DBClient.HGet(shortDB.DBContext, surl, "longurl").Result()
	if err != nil {
		return false, err
	}
	return rec == lurl, err
}

func (shortDB ShortDB) GetCreationTimestamp(surl string) (string, error) {
	return shortDB.DBClient.HGet(shortDB.DBContext, surl, "creation").Result()
}

func (shortDB ShortDB) GetLastAccess(surl string) (string, error) {
	return shortDB.DBClient.HGet(shortDB.DBContext, surl, "last").Result()
}

func (shortDB ShortDB) GetAccessCount(surl string) (string, error) {
	return shortDB.DBClient.HGet(shortDB.DBContext, surl, "count").Result()
}
