package main

import (
	"fmt"
	"log"
	"net/http"

	app "shorturl/internal/app"
	web "shorturl/internal/web"
)

func main() {
	app.Config.FromEnv()

	app.DBInstance.Init()

	pong, err := app.DBInstance.Ping()
	fmt.Println(pong, err)

	http.HandleFunc("/", web.UIHandler)
	log.Fatal(http.ListenAndServe(":"+app.Config.AppPort, nil))
}
