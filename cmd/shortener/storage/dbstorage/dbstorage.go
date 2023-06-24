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

func (s *DBStorageStruct) PutURL(id string, body string, userID interface{}) (string, error) {

	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return "", err
	}
	defer conn.Close(ctx)

	// TODO Артем https://github.com/GoogleCloudPlatform/pgadapter/blob/postgresql-dialect/docs/pgx.md#batching

	var _id string
	err = conn.QueryRow(ctx, "INSERT INTO public.short_url (id, original_url, created_by) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO UPDATE SET original_url=EXCLUDED.original_url RETURNING id;", id, body, userID).Scan(&_id)

	if err != nil {
		return id, err
	}

	if id != _id {
		err = ErrConflict
		return _id, err
	}

	return id, err

}

func (s *DBStorageStruct) PutURLBatch(body []models.RequestBatch, userID interface{}) ([]models.ResponseButch, error) {
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
	}(conn, ctx)

	batch := &pgx.Batch{}
	for i := 0; i < len(body); i++ {
		batch.Queue("INSERT INTO short_url (id, original_url, created_by) VALUES($1, $2, $3)", body[i].UUID, body[i].URL, userID)
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

func (s *DBStorageStruct) GetURLBatch(userID interface{}) ([]models.ResponseButchForUser, error) {
	bodyResponseButch := make([]models.ResponseButchForUser, 0)
	fmt.Println(bodyResponseButch)
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
	}(conn, ctx)

	rows, err := conn.Query(ctx, "select id, original_url from public.short_url where created_by=$1", userID) //.Scan(&ID, &OriginalURL)
	if err != nil {
		sugar.Errorf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var ID string
		var OriginalURL string
		err = rows.Scan(&ID, &OriginalURL)
		if err != nil {
			return nil, err
		}

		row := models.ResponseButchForUser{}

		row.ShortURL = fmt.Sprintf("%s/%s", cfg.FlagBaseURLAddr, ID)
		row.OriginalURL = OriginalURL

		fmt.Println(ID, OriginalURL)

		bodyResponseButch = append(bodyResponseButch, row)
	}
	//for i := 0; i < len(body); i++ {
	//	bodyResponseButch[i].UUID = body[i].UUID
	//	bodyResponseButch[i].URL = fmt.Sprintf("%s/%s", cfg.FlagBaseURLAddr, body[i].UUID)
	//}

	return bodyResponseButch, nil
	//return nil, nil
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
	defer conn.Close(ctx)

	var OriginalURL string
	err = conn.QueryRow(ctx, "select original_url from short_url where id=$1 and is_deleted=false", id).Scan(&OriginalURL)
	if err != nil {
		sugar.Errorf("QueryRow failed: %v\n", err)
		return "", err
	}

	return OriginalURL, nil
}

func (s *DBStorageStruct) InitStorage() error {
	return nil
}

func (s *DBStorageStruct) DeleteURLBatch(req []string, userID interface{}) error {
	sugar := logging.GetSugarLogger()
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.FlagDataBaseDSN)
	if err != nil {
		sugar.Errorf("Unable to connect to database: %v\n", err)
		return err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			sugar.Errorf("Error %s", err)
		}
	}(conn, ctx)

	batch := &pgx.Batch{}
	for i := 0; i < len(req); i++ {
		batch.Queue("UPDATE public.short_url\nSET is_deleted = false\nWHERE id = $1 and created_by = $2;", req[i], userID)
	}
	br := conn.SendBatch(ctx, batch)
	_, err = br.Exec()
	//if err != nil {
	//	return err
	//}
	return err
}
