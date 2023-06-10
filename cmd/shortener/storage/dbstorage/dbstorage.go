package dbstorage

import (
	"context"

	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorageStruct struct{}

func (s *DBStorageStruct) PutURL(body string) (string, error) {
	id := uuid.New().String()

	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return "", err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO public.short_url (id, original_url)VALUES ($1, $2);", id, body)
	if err != nil {
		sugar.Errorf("Createuuid extension failed: %v\n", err)
		return "", err
	}

	return id, nil
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
	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	//var name string
	// Устанавливаем расширение для uuid.
	_, err = conn.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		sugar.Errorf("Createuuid extension failed: %v\n", err)
		return err
	}

	// Создаем таблицу для хранения.
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS short_url (id UUID PRIMARY KEY, original_url TEXT NOT NULL);")
	if err != nil {
		sugar.Errorf("Create table failed %v\n", err)
		return err
	}
	return err
}
