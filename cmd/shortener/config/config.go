package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr     string
	FlagBaseURLAddr string
}

//var config = &configuration{}

func NewConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.FlagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&cfg.FlagBaseURLAddr, "b", "http://127.0.0.1:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.FlagRunAddr = envRunAddr
	}

	if envBaseURLAddr := os.Getenv("BASE_URL"); envBaseURLAddr != "" {
		cfg.FlagBaseURLAddr = envBaseURLAddr
	}

	return cfg
}
