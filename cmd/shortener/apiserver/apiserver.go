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

type localStorage struct{}

func (s *localStorage) PutURL(body string) string {
	id := uuid.New().String()
	shortURL[id] = body
	return id
}

func (s *localStorage) GetURL(id string) string {
	return shortURL[id]
}

var (
	shortURL = map[string]string{
		"094c4130-9674-4c18-bf60-7385d7f61934": "https://practicum.yandex.ru/",
	}
)

func createShortURL(rw http.ResponseWriter, r *http.Request) {

	log.Println("Create short url")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ls storage.AbstractStorage = &localStorage{}
	id := ls.PutURL(string(bodyBytes))

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprintf(rw, "%s/%s", config.Config.FlagBaseURLAddr, id)
	if err != nil {
		return
	}
}

func readShortURL(rw http.ResponseWriter, r *http.Request) {
	var ls storage.AbstractStorage = &localStorage{}

	rw.Header().Set("Content-Type", "text/plain")
	log.Printf("Get short url %s xxx", ls.GetURL(chi.URLParam(r, "id")))

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
