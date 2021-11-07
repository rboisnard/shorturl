package app

import "os"

type AppConfig struct {
	AppURL    string
	AppIP     string
	AppPort   string
	RedisIP   string
	RedisPort string
}

// default values
var Config = AppConfig{
	AppURL:    "localhost:8080",
	AppIP:     "localhost",
	AppPort:   "5500",
	RedisIP:   "localhost",
	RedisPort: "6379",
}

func (config *AppConfig) FromEnv() {
	AppURL, found := os.LookupEnv("APP_URL")
	if found {
		config.AppURL = AppURL
	}

	AppIP, found := os.LookupEnv("APP_IP")
	if found {
		config.AppIP = AppIP
	}

	AppPort, found := os.LookupEnv("APP_PORT")
	if found {
		config.AppPort = AppPort
	}

	RedisIP, found := os.LookupEnv("REDIS_IP")
	if found {
		config.RedisIP = RedisIP
	}

	RedisPort, found := os.LookupEnv("REDIS_PORT")
	if found {
		config.RedisPort = RedisPort
	}
}
