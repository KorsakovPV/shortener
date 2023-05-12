package main

import (
	"fmt"
	"github.com/KorsakovPV/shortener/cmd/shortener/config"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type abstractStorage interface {
	putURL(string) string
	getURL(string) string
}

type localStorage struct{}

func (s *localStorage) putURL(body string) string {
	id := uuid.New().String()
	shortURL[id] = body
	return id
}

func (s *localStorage) getURL(id string) string {
	return shortURL[id]
}

var (
	shortURL = map[string]string{
		"094c4130-9674-4c18-bf60-7385d7f61934": "https://practicum.yandex.ru/",
	}
)

func createShortURL(rw http.ResponseWriter, r *http.Request) {
	//cfg := config.NewConfig()

	log.Println("Create short url")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ls abstractStorage = &localStorage{}
	id := ls.putURL(string(bodyBytes))

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, _ = rw.Write([]byte(fmt.Sprintf("%s/%s", config.Config.FlagBaseURLAddr, id)))
	//_, err = fmt.Fprintf(rw, "%s/%s", cfg.FlagBaseURLAddr, id)
	//if err != nil {
	//	return
	//}
}

func readShortURL(rw http.ResponseWriter, r *http.Request) {
	log.Println("Get short url")

	var ls abstractStorage = &localStorage{}

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Location", ls.getURL(chi.URLParam(r, "id")))
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func methodNotAllowed(rw http.ResponseWriter, r *http.Request) {
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

func main() {

	config.ParseFlags()

	log.Printf("Shortener start on %s. Default base URL %s.", config.Config.FlagRunAddr, config.Config.FlagBaseURLAddr)

	log.Fatal(http.ListenAndServe(config.Config.FlagRunAddr, Router()))
}
