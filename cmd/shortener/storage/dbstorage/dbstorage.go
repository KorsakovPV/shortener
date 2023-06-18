package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/internal/models"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrConflict = errors.New("data conflict")

type DBStorageStruct struct{}

func (s *DBStorageStruct) PutURL(id string, body string) (string, error) {

	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return "", err
	}
	defer conn.Close(ctx)

	var _id string
	err = conn.QueryRow(ctx, "INSERT INTO public.short_url (id, original_url) VALUES ($1, $2) ON CONFLICT (original_url) DO UPDATE SET original_url=EXCLUDED.original_url RETURNING id;", id, body).Scan(&_id)

	if err != nil {
		return id, err
	}

	if id != _id {
		err = ErrConflict
		return _id, err
	}

	return id, err

}

func (s *DBStorageStruct) PutURLBatch(body []models.RequestBatch) ([]models.ResponseButch, error) {
	bodyResponseButch := make([]models.ResponseButch, len(body))
	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			sugar.Errorf("Error %s", err)
		}
	}(conn, context.Background())

	batch := &pgx.Batch{}
	for i := 0; i < len(body); i++ {
		batch.Queue("INSERT INTO short_url (id, original_url) VALUES($1, $2)", body[i].UUID, body[i].URL)
	}
	br := conn.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(body); i++ {
		bodyResponseButch[i].UUID = body[i].UUID
		bodyResponseButch[i].URL = fmt.Sprintf("%s/%s", cfg.FlagBaseURLAddr, body[i].UUID)
	}

	return bodyResponseButch, nil
}

func (s *DBStorageStruct) GetURL(id string) (string, error) {
	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return "", err
	}
	defer conn.Close(context.Background())

	var OriginalURL string
	err = conn.QueryRow(context.Background(), "select original_url from short_url where id=$1", id).Scan(&OriginalURL)
	if err != nil {
		sugar.Errorf("QueryRow failed: %v\n", err)
		return "", err
	}

	return OriginalURL, nil
}

func (s *DBStorageStruct) InitStorage() error {
	return nil
	//sugar := logging.GetSugarLogger()
	//cfg := config.GetConfig()
	//ctx := context.Background()
	//
	//conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	//if err != nil {
	//	sugar.Errorf("Unable to connect to database: %v\n", err)
	//	return err
	//}
	//defer conn.Close(context.Background())
	//
	//// Устанавливаем расширение для uuid.
	//_, err = conn.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	//if err != nil {
	//	sugar.Errorf("Createuuid extension failed: %v\n", err)
	//	return err
	//}
	//
	//// Создаем таблицу для хранения.
	//_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS short_url (id TEXT PRIMARY KEY, original_url TEXT NOT NULL UNIQUE);")
	//if err != nil {
	//	sugar.Errorf("Create table failed %v\n", err)
	//	return err
	//}
	//return err
}
