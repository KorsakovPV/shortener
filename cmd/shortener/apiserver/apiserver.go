package apiserver

import (
	"fmt"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type LocalStorage struct {
	shortURL map[string]string
}

func (s *LocalStorage) PutURL(body string) string {
	id := uuid.New().String()
	s.shortURL[id] = body
	return id
}

func (s *LocalStorage) GetURL(id string) string {
	return s.shortURL[id]
}

func createShortURL(rw http.ResponseWriter, r *http.Request) {

	log.Println("Create short url")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ls storage.AbstractStorage = &LocalStorage{}
	id := ls.PutURL(string(bodyBytes))

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprintf(rw, "%s/%s", config.Config.FlagBaseURLAddr, id)
	if err != nil {
		return
	}
}

func readShortURL(rw http.ResponseWriter, r *http.Request) {
	log.Println("Get short url")

	var ls storage.AbstractStorage = &LocalStorage{}

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Location", ls.GetURL(chi.URLParam(r, "id")))
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