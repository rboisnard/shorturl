package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	app "shorturl/internal/app"
	web "shorturl/internal/web"

	"github.com/go-redis/redis/v8"
)

func main() {
	app.Config.SetAppConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     app.Config.RedisIP + ":" + app.Config.RedisPort,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(context.TODO()).Result()
	fmt.Println(pong, err)

	http.HandleFunc("/", web.UIHandler)
	log.Fatal(http.ListenAndServe(":"+app.Config.AppPort, nil))
}
