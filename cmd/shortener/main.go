package main

import "net/http"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, http.HandlerFunc(webhook))

	//return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
	return http.ListenAndServe(`:8080`, mux)
}

func webhook(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTemporaryRedirect)

			// пока установим ответ-заглушку, без проверки ошибок
			_, _ = w.Write([]byte(`
      {
        "response": {
          "Location": "https://practicum.yandex.ru/"
        },
        "version": "1.0"
      }
    `))
		}

	case http.MethodPost:
		{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			// пока установим ответ-заглушку, без проверки ошибок
			_, _ = w.Write([]byte(`
      {
        "response": {
          "url": "http://localhost:8080/EwHXdJfB"
        },
        "version": "1.0"
      }
    `))
		}

	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
