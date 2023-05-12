package main

import (
	api "github.com/KorsakovPV/shortener/cmd/shortener/api"
	config "github.com/KorsakovPV/shortener/cmd/shortener/config"
	//"github.com/xlab/closer"
	"log"
	"net/http"
)

func main() {
	//closer.Bind(cleanup)
	//_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//defer stop()

	cfg := config.NewConfig()

	log.Printf("Shortener start on %s. Default base URL %s.", cfg.FlagRunAddr, cfg.FlagBaseURLAddr)

	log.Fatal(http.ListenAndServe(cfg.FlagRunAddr, api.Router()))
}
