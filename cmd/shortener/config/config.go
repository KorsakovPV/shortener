package config

import (
	"flag"
	"os"
)

type Сonfiguration struct {
	FlagRunAddr         string
	FlagBaseURLAddr     string
	FlagFileStoragePath string
	//FlagDataBaseDSN     string
}

var config = &Сonfiguration{}

func GetConfig() *Сonfiguration {
	return config
}

func ParseFlags() {
	flag.StringVar(&config.FlagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&config.FlagBaseURLAddr, "b", "http://127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&config.FlagFileStoragePath, "f", "/tmp/short-url-db.json", "address and port to run server")
	//flag.StringVar(&config.FlagDataBaseDSN, "d", "", "address and port to run server")
	//flag.StringVar(&config.FlagDataBaseDSN, "d", "postgres://postgres:postgres@localhost:5432/shortener", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.FlagRunAddr = envRunAddr
	}

	if envBaseURLAddr := os.Getenv("BASE_URL"); envBaseURLAddr != "" {
		config.FlagBaseURLAddr = envBaseURLAddr
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		config.FlagFileStoragePath = envFileStoragePath
	}

	//if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
	//	config.FlagDataBaseDSN = envDataBaseDSN
	//}
}
