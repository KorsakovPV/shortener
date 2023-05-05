package main

import (
	"flag"
)

type configuration struct {
	flagRunAddr     string
	flagBaseURLAddr string
}

var config = &configuration{}

func parseFlags() {
	flag.StringVar(&config.flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.flagBaseURLAddr, "b", "localhost:8080", "address and port to run server")
	flag.Parse()
}
