package main

import (
	"github.com/KorsakovPV/shortener/cmd/shortener/apiserver"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
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
	)
	err := http.ListenAndServe(cfg.FlagRunAddr, apiserver.Router())
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

}
