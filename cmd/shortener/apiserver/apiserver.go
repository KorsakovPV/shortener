package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/logging"
	"github.com/KorsakovPV/shortener/cmd/shortener/middleware"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/KorsakovPV/shortener/internal/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
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

		id := storage.GetStorage().PutURL(string(bodyBytes))

		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusCreated)

		_, err = fmt.Fprintf(rw, "%s/%s", config.GetConfig().FlagBaseURLAddr, id)
		if err != nil {
			sugar.Errorf("ERROR Can't writing content to HTTP response. %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	return http.HandlerFunc(fn)
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

		id := storage.GetStorage().PutURL(req.Url)

		resp := models.Response{
			Result: fmt.Sprintf("%s/%s", config.GetConfig().FlagBaseURLAddr, id),
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			sugar.Debug("error encoding response", zap.Error(err))
			return
		}

	}
	return http.HandlerFunc(fn)
}

func readShortURL() http.HandlerFunc {
	sugar := logging.GetSugarLogger()

	fn := func(rw http.ResponseWriter, r *http.Request) {
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
	sugar := logging.GetSugarLogger()

	fn := func(rw http.ResponseWriter, _ *http.Request) {
		sugar.Errorln("Method Not Allowed")
		rw.WriteHeader(http.StatusBadRequest)
	}
	return http.HandlerFunc(fn)
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", middleware.WithLogging(createShortURL()))
	r.Post("/api/shorten", middleware.WithLogging(createShortURLJson()))
	r.Get("/{id}", middleware.WithLogging(readShortURL()))
	r.MethodNotAllowed(middleware.WithLogging(methodNotAllowed()))
	return r
}
