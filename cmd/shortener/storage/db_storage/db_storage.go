package db_storage

import (
	"context"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDB(ctx context.Context) error {
	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()

	conn, err := pgx.Connect(context.Background(), cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return err
	}
	err = conn.Ping(ctx)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return err
	}
	return nil
}
