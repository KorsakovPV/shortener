package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/apiserver"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"log"
	"net/http"
)

//type AbstractStorage interface {
//	PutURL(string) string
//	GetURL(string) string
//}

func main() {

	config.ParseFlags()

	log.Printf("Shortener start on %s. Default base URL %s.", config.Config.FlagRunAddr, config.Config.FlagBaseURLAddr)

	log.Fatal(http.ListenAndServe(config.Config.FlagRunAddr, apiserver.Router()))
}
