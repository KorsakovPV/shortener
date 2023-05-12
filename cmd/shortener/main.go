package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/api"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"log"
	"net/http"
)

func main() {

	cfg := config.NewConfig()

	log.Printf("Shortener start on %s. Default base URL %s.", cfg.FlagRunAddr, cfg.FlagBaseURLAddr)

	log.Fatal(http.ListenAndServe(cfg.FlagRunAddr, api.Router()))
}
