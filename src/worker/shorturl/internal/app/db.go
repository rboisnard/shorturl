package app

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type ShortDB struct {
	DBClient  *redis.Client
	DBContext context.Context
	maxRetry  int
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
	shortDB.maxRetry = 3
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

func (shortDB ShortDB) StoreLongURL(lurl string) (string, error) {
	// internal errors

	var didNotExistError = errors.New("Short URL did not exist")
	var foundSimilarError = errors.New("Similar short URL already exists")
	var foundDifferentError = errors.New("Short URL already exists for a different long URL")
	var noHashAvailableError = errors.New("No short URL available for the long URL")

	// StoreLongURL logic starts here
	abort := make(chan struct{})
	results := Shorten(abort, lurl)

	for surl := range results {
		log.Printf("+ testing storage for short URL '%s'", surl)

		// internal function for transactions
		addLongURL := func(tx *redis.Tx) error {
			res, err := tx.HGet(shortDB.DBContext, surl, "longurl").Result()
			if err == redis.Nil {
				log.Printf("++ no record found")
				// no record was found, ok to try storing now in
				// an internal function, pipeline is using multi/exec
				_, err := tx.TxPipelined(shortDB.DBContext, func(pipe redis.Pipeliner) error {
					var now string = strconv.FormatInt(time.Now().Unix(), 10)

					_, err := pipe.HSet(shortDB.DBContext, surl, []string{
						"longurl", lurl,
						"creation", now,
						"last", now,
						"count", "0",
					}).Result()

					return err
				})

				if err != nil {
					log.Printf("++ tx error '%s'", err.Error())
					return err
				}
				log.Printf("++ tx ok")
				return didNotExistError
			}

			if err == nil {
				// a record for surl was found
				if res == lurl {
					log.Printf("+ record found, similar '%s'", res)
					return foundSimilarError
				} else {
					log.Printf("+ record found, different '%s'", res)
					return foundDifferentError
				}
			}

			// err is not nil, there was some issue
			log.Printf("+ unknown issue '%s'", err.Error())
			return err
		}

		// internal function to check issues and retry
		hadIssues := func(err error) bool {
			switch err {
			case didNotExistError, foundSimilarError, foundDifferentError:
				return false
			default:
				return true
			}
		}

		// start a watch on a multi/exec transaction
		currentTrial, err := 0, shortDB.DBClient.Watch(shortDB.DBContext, addLongURL, surl)

		for currentTrial < shortDB.maxRetry && hadIssues(err) {
			currentTrial++
			log.Printf("+ trying again (#%d)", currentTrial)
			err = shortDB.DBClient.Watch(shortDB.DBContext, addLongURL, surl)
		}

		switch err {
		case didNotExistError:
			// store was succesful
			close(abort)
			log.Printf("+ success with '%s'", surl)
			return surl, error(nil)

		case foundSimilarError:
			// this long URL already exists
			// update the timestamp and return
			// the short URL
			close(abort)
			_, err = shortDB.GetLongURL(surl)
			log.Printf("+ success with '%s' (existed)", surl)
			return surl, err

		case foundDifferentError:
			// this short URL already exists
			// try another one
			log.Printf("+ try a new short URL, '%s' already exists", surl)
			continue

		default:
			// something happened
			log.Printf("+ unknown issue '%s'", err.Error())
			return surl, err
		}
	}
	log.Printf("+ no hash available for '%s'", lurl)
	return "error", noHashAvailableError
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
