package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handler)

	http.ListenAndServe(":8000", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hoge"))
	w.WriteHeader(http.StatusOK)
	return
}
