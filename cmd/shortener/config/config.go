package config

import (
	"flag"
	"os"
)

type Сonfiguration struct {
	FlagRunAddr     string
	FlagBaseURLAddr string
}

var Config = &Сonfiguration{}

//func init() {
//	ParseFlags()
//}

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

func GetConfig() *Сonfiguration {
	return Config
}
