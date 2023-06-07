package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/KorsakovPV/shortener/internal/apiserver"
	"net/http"
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
	)

	if cfg.FlagFileStoragePath != "" {
		err := storage.GetStorage().LoadBackupURL()
		if err != nil {
			sugar.Errorf("ERROR LoadBackupURL. %s", err)
			return
		}
	}

	err := http.ListenAndServe(cfg.FlagRunAddr, apiserver.Router())

	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

}
