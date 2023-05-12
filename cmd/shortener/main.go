package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/api"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"log"
	"net/http"
)

func main() {

	config.ParseFlags()

	log.Printf("Shortener start on %s. Default base URL %s.", config.Config.FlagRunAddr, config.Config.FlagBaseURLAddr)

	log.Fatal(http.ListenAndServe(config.Config.FlagRunAddr, api.Router()))
}
