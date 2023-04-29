package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, http.HandlerFunc(webhook))

	return http.ListenAndServe(`:8080`, mux)
}

func webhook(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPost:
		{
			log.Println("Create short url")
			id := uuid.New().String()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			shortURL[id] = string(bodyBytes)

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)

			_, _ = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))
		}

	case http.MethodGet:
		{
			log.Println("Get short url")

			id := strings.Split(r.RequestURI, "/")[len(strings.Split(r.RequestURI, "/"))-1]

			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Location", shortURL[id])
			w.WriteHeader(http.StatusTemporaryRedirect)

		}

	default:
		{
			log.Println("Bad Request")

			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
