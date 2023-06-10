package main

import (
	"net/http"

	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/KorsakovPV/shortener/internal/apiserver"
)

func main() {
	sugar := logging.GetSugarLogger()

	config.ParseFlags()
	cfg := config.GetConfig()
	sugar.Infow(
		"Starting server",
		"address", cfg.FlagRunAddr,
		"Default base URL", cfg.FlagBaseURLAddr,
		"File for store", cfg.FlagFileStoragePath,
		"DataBase DSN", cfg.FlagDataBaseDSN,
	)

	err := storage.InitStorage()
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "init storage")
	}

	err = http.ListenAndServe(cfg.FlagRunAddr, apiserver.Router())

	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

}
