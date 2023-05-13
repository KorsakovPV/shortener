package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/apiserver"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"log"
	"net/http"
)

func main() {

	config.ParseFlags()

	log.Printf("Shortener start on %s. Default base URL %s.", config.Config.FlagRunAddr, config.Config.FlagBaseURLAddr)

	err := http.ListenAndServe(config.Config.FlagRunAddr, apiserver.Router())
	if err != nil {
		log.Fatal(err)
	}

}
