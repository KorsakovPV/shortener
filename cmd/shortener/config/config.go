package config

import (
	"flag"
	"os"
)

type configuration struct {
	FlagRunAddr     string
	FlagBaseURLAddr string
}

var Config = &configuration{}

func ParseFlags() {
	flag.StringVar(&Config.FlagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&Config.FlagBaseURLAddr, "b", "http://127.0.0.1:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Config.FlagRunAddr = envRunAddr
	}

	if envBaseURLAddr := os.Getenv("BASE_URL"); envBaseURLAddr != "" {
		Config.FlagBaseURLAddr = envBaseURLAddr
	}
}
