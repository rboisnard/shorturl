package main

import (
	"fmt"
	"log"
	"net/http"

	app "shorturl/internal/app"
	web "shorturl/internal/web"
)

func main() {
	app.Config.SetAppConfig()

	var db app.ShortDB
	db.Init()

	pong, err := db.Ping()
	fmt.Println(pong, err)

	http.HandleFunc("/", web.UIHandler)
	log.Fatal(http.ListenAndServe(":"+app.Config.AppPort, nil))
}
