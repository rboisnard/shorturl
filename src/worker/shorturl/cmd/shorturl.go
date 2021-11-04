package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

type AppConfig struct {
	appURL    string
	appIP     string
	appPort   string
	redisIP   string
	redisPort string
}

// default values
var appConfig = AppConfig{
	appURL:    "localhost:8080",
	appIP:     "localhost",
	appPort:   "5500",
	redisIP:   "localhost",
	redisPort: "6379",
}

type HomeData struct {
	HasError     bool
	HasResult    bool
	HasForm      bool
	BaseURL      string
	ShortURL     string
	ErrorMessage string
}

func main() {
	setAppConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.redisIP + ":" + appConfig.redisPort,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(context.TODO()).Result()
	fmt.Println(pong, err)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+appConfig.appPort, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// common handler for all requests

	if len(r.URL.Path) == 1 {

		// default requested path is "/"
		if r.Method == "GET" {
			homeGET(w, r)
		} else if r.Method == "POST" {
			homePOST(w, r)
		} else {
			internalServerError(w, r, "unexpected method '"+r.Method+"' for '/'")
		}
		return

	} else if len(r.URL.Path) > 1 {

		// requested path contains a shorturl
		shorturl := r.URL.Path[1:]

		if r.Method == "GET" {
			shorturlGET(w, r, shorturl)
		} else {
			internalServerError(w, r, "unexpected method '"+r.Method+"' for '/"+shorturl+"'")
		}
		return

	}

	// shouldn't be here
	internalServerError(w, r, "reached handler end without returning first")
}

func homeGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("request GET '/' from %s", r.RemoteAddr)

	data := HomeData{
		HasError:     false,
		HasResult:    false,
		HasForm:      true,
		BaseURL:      "",
		ShortURL:     "",
		ErrorMessage: "",
	}
	render(w, r, data)

	log.Printf("response to GET '/' from %s", r.RemoteAddr)
}

func homePOST(w http.ResponseWriter, r *http.Request) {
	log.Printf("request POST '/' from %s with '%s'", r.RemoteAddr, r.FormValue("longurl"))

	// TODO: fetch/store the shortened url
	// tmp: mock shorturl result
	shorturl, err := "mock_url", error(nil)

	if err != nil {
		// error with the input long URL
		data := HomeData{
			HasError:     true,
			HasResult:    false,
			HasForm:      true,
			BaseURL:      "",
			ShortURL:     "",
			ErrorMessage: "the input long URL resulted in error",
		}
		render(w, r, data)

		log.Printf("response to POST '/' from %s - bad long URL", r.RemoteAddr)
		return
	}

	// no error, reply with the shortened URL
	data := HomeData{
		HasError:     false,
		HasResult:    true,
		HasForm:      false,
		BaseURL:      "http://" + appConfig.appURL,
		ShortURL:     shorturl,
		ErrorMessage: "",
	}
	render(w, r, data)

	log.Printf("response to POST '/' from %s - shorturl '%s'", r.RemoteAddr, shorturl)
}

func shorturlGET(w http.ResponseWriter, r *http.Request, shorturl string) {
	log.Printf("request GET '/%s' from %s", shorturl, r.RemoteAddr)

	// TODO: fetch and validate the shorturl
	// tmp: mock shorturl result
	longurl, err := "http://www.google.com", error(nil)

	if err != nil {
		// error with the input long URL
		data := HomeData{
			HasError:     true,
			HasResult:    false,
			HasForm:      true,
			BaseURL:      "",
			ShortURL:     "",
			ErrorMessage: "the input short URL resulted in error",
		}
		render(w, r, data)

		log.Printf("response to GET '/%s' from %s - bad short URL", shorturl, r.RemoteAddr)
		return
	}

	// no error, simply redirect
	log.Printf("response to GET '/%s' from %s - redirect to '%s'", shorturl, r.RemoteAddr, longurl)
	http.Redirect(w, r, longurl, http.StatusFound)
}

func render(w http.ResponseWriter, r *http.Request, data HomeData) {
	home, err := template.ParseFiles("home.html")
	if err != nil {
		internalServerError(w, r, "cannot open home.html")
		return
	}

	err = home.Execute(w, data)
	if err != nil {
		internalServerError(w, r, "error while rendering home.html")
		return
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, reason string) {
	var request bytes.Buffer
	r.Write(&request)
	log.Printf("error '%s' with '%s'", reason, request.String())
	http.Error(w, "unexpected behavior", http.StatusInternalServerError)
}

func setAppConfig() {
	appURL, found := os.LookupEnv("APP_URL")
	if found {
		appConfig.appURL = appURL
	}

	appIP, found := os.LookupEnv("APP_IP")
	if found {
		appConfig.appIP = appIP
	}

	appPort, found := os.LookupEnv("APP_PORT")
	if found {
		appConfig.appPort = appPort
	}

	redisIP, found := os.LookupEnv("REDIS_IP")
	if found {
		appConfig.redisIP = redisIP
	}

	redisPort, found := os.LookupEnv("REDIS_PORT")
	if found {
		appConfig.redisPort = redisPort
	}
}
