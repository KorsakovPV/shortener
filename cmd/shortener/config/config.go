package config

import (
	"flag"
	"os"
)

type Сonfiguration struct {
	FlagRunAddr     string
	FlagBaseURLAddr string
}

var config = &Сonfiguration{}

func init() {
	ParseFlags()
}

func ParseFlags() {
	flag.StringVar(&config.FlagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&config.FlagBaseURLAddr, "b", "http://127.0.0.1:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.FlagRunAddr = envRunAddr
	}

	if envBaseURLAddr := os.Getenv("BASE_URL"); envBaseURLAddr != "" {
		config.FlagBaseURLAddr = envBaseURLAddr
	}
}

func GetConfig() *Сonfiguration {
	return config
}
