package apiserver

import (
	"fmt"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

func createShortURL(rw http.ResponseWriter, r *http.Request) {
	log.Println("Create short url")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	id := storage.LocalStorage.PutURL(string(bodyBytes))

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprintf(rw, "%s/%s", config.Config.FlagBaseURLAddr, id)
	if err != nil {
		return
	}
}

func readShortURL(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	log.Printf("Get short url %s xxx", storage.LocalStorage.GetURL(chi.URLParam(r, "id")))

	rw.Header().Set("Location", storage.LocalStorage.GetURL(chi.URLParam(r, "id")))
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func methodNotAllowed(rw http.ResponseWriter, _ *http.Request) {
	log.Println("Method Not Allowed")
	rw.WriteHeader(http.StatusBadRequest)
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", readShortURL)
	r.Post("/", createShortURL)
	r.MethodNotAllowed(methodNotAllowed)
	return r
}
