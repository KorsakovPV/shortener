package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/apiserver"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"log"
	"net/http"
)

func main() {

	cfg := config.GetConfig()

	log.Printf("Shortener start on %s. Default base URL %s.", cfg.FlagRunAddr, cfg.FlagBaseURLAddr)

	err := http.ListenAndServe(cfg.FlagRunAddr, apiserver.Router())
	if err != nil {
		log.Fatal(err)
	}

}
