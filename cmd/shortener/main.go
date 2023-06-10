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
		//"DataBase DSN", cfg.FlagDataBaseDSN,
	)

	err := storage.InitStorage()
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "init storage")
	}

	//dbstorage.Connect()

	//// urlExample := "postgres://username:password@localhost:5432/database_name"
	//conn, err := pgx.Connect(context.Background(), cfg.FlagDataBaseDSN)
	//if err != nil {
	//	sugar.Errorf("Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	//
	////var name string
	////var weight int64
	////var table_catalog string
	////var table_name string
	//var number int
	//err = conn.QueryRow(context.Background(), "select 0").Scan(&number) //, &table_name)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//fmt.Println(number)

	// TODO если переданы параметры БД метод не должен отрабатывать.
	//if cfg.FlagFileStoragePath != "" {
	//	err := storage.GetStorage().LoadBackupURL()
	//	if err != nil {
	//		sugar.Errorf("ERROR LoadBackupURL. %s", err)
	//		return
	//	}
	//}

	err = http.ListenAndServe(cfg.FlagRunAddr, apiserver.Router())

	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

}
