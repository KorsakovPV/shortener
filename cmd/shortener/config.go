package main

import (
	"flag"
	"os"
)

type configuration struct {
	flagRunAddr     string
	flagBaseURLAddr string
}

var config = &configuration{}

func parseFlags() {
	flag.StringVar(&config.flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.flagBaseURLAddr, "b", "http://localhost:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		config.flagRunAddr = envRunAddr
	}

	if envBaseURLAddr := os.Getenv("RUN_ADDR"); envBaseURLAddr != "" {
		config.flagBaseURLAddr = envBaseURLAddr
	}
}
