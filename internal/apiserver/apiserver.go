package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/cmd/shortener/middleware"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage/dbstorage"
	"github.com/KorsakovPV/shortener/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func createShortURL() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		sugar := logging.GetSugarLogger()

		sugar.Infoln("Create short url")

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			sugar.Errorf("ERROR Can't get value from body. %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		id := uuid.New().String()

		id, err = storage.GetStorage().PutURL(id, string(bodyBytes))

		switch {
		case errors.Is(err, dbstorage.ErrConflict):
			{
				rw.Header().Set("Content-Type", "text/plain")
				rw.WriteHeader(http.StatusConflict)
			}
		case err != nil:
			{
				sugar.Errorf("ERROR Can't writing content to HTTP response. %s", err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			{
				rw.Header().Set("Content-Type", "text/plain")
				rw.WriteHeader(http.StatusCreated)
			}
		}

		_, err = fmt.Fprintf(rw, "%s/%s", config.GetConfig().FlagBaseURLAddr, id)
		if err != nil {
			sugar.Errorf("ERROR Can't writing content to HTTP response. %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func pingDB() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		sugar := logging.GetSugarLogger()

		sugar.Infoln("Ping DB.")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := PingDB(ctx)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(http.StatusOK)
		}

	}
	return http.HandlerFunc(fn)
}

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

func createShortURLJson() http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
		sugar := logging.GetSugarLogger()

		sugar.Infoln("Create short url")

		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			sugar.Debug("cannot decode request JSON body", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		sugar.Infoln(req)

		id := uuid.New().String()

		id, err := storage.GetStorage().PutURL(id, req.URL)

		resp := models.Response{
			Result: fmt.Sprintf("%s/%s", config.GetConfig().FlagBaseURLAddr, id),
		}

		switch {
		case errors.Is(err, dbstorage.ErrConflict):
			{
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusConflict)
			}
		case err != nil:
			{
				sugar.Errorf("ERROR Can't writing content to HTTP response. %s", err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			{
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusCreated)
			}
		}

		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			sugar.Debug("error encoding response", zap.Error(err))
			return
		}
		sugar.Infoln(resp)
	}
	return http.HandlerFunc(fn)
}

func createShortURLBatchJSON() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		sugar := logging.GetSugarLogger()

		sugar.Infoln("Create batch short url")

		var req []models.RequestBatch
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			sugar.Debug("cannot decode request JSON body", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		sugar.Infoln(req)

		bodyResponseButch, err := storage.GetStorage().PutURLBatch(req)
		if err != nil {
			sugar.Errorf("ERROR Can't writing content to HTTP response. %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		resp := bodyResponseButch

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			sugar.Debug("error encoding response", zap.Error(err))
			return
		}
		sugar.Infoln(resp)
	}
	return http.HandlerFunc(fn)
}

func readShortURL() http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
		sugar := logging.GetSugarLogger()

		sugar.Infof("Get URL. id=%s", chi.URLParam(r, "id"))

		shortURL, err := storage.GetStorage().GetURL(chi.URLParam(r, "id"))

		if err != nil {
			sugar.Errorf("ERROR %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		rw.Header().Set("Content-Type", "text/plain")
		sugar.Infof("Get short url %s", shortURL)

		rw.Header().Set("Location", shortURL)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
	return http.HandlerFunc(fn)
}

func methodNotAllowed() http.HandlerFunc {

	fn := func(rw http.ResponseWriter, _ *http.Request) {
		sugar := logging.GetSugarLogger()

		sugar.Errorln("Method Not Allowed")
		rw.WriteHeader(http.StatusBadRequest)
	}
	return http.HandlerFunc(fn)
}

func middlewares(h http.HandlerFunc) http.HandlerFunc {
	return middleware.WithLogging(middleware.GzipMiddleware(h))
}

func Router() chi.Router {

	r := chi.NewRouter()

	r.Post("/api/shorten", middlewares(createShortURLJson()))
	r.Post("/api/shorten/batch", middlewares(createShortURLBatchJSON()))
	r.Get("/ping", middlewares(pingDB()))
	r.Get("/{id}", middlewares(readShortURL()))
	r.Post("/", middlewares(createShortURL()))
	r.MethodNotAllowed(middlewares(methodNotAllowed()))

	return r
}
