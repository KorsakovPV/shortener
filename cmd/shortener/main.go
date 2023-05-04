package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

var (
	shortURL = map[string]string{
		"094c4130-9674-4c18-bf60-7385d7f61934": "https://practicum.yandex.ru/",
	}
)

func createShortURL(rw http.ResponseWriter, r *http.Request) {
	log.Println("Create short url")
	id := uuid.New().String()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	shortURL[id] = string(bodyBytes)

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, _ = rw.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))
}

func readShortURL(rw http.ResponseWriter, r *http.Request) {
	log.Println("Get short url")

	//id := strings.Split(r.RequestURI, "/")[len(strings.Split(r.RequestURI, "/"))-1]
	id := chi.URLParam(r, "id")

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Location", shortURL[id])
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

	log.Fatal(http.ListenAndServe(":8080", Router()))
}
