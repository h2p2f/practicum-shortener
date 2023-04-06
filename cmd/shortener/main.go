package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-shortener/internal/app/handler"
	"github.com/h2p2f/practicum-shortener/internal/app/storage"
	"log"
	"net/http"
)

func shortenerRouter() chi.Router {
	s := storage.NewLinkStorage()
	handlers := handler.NewStorageHandler(s)
	r := chi.NewRouter()
	r.Post("/", handlers.PostLinkHandler)
	r.Get("/{id}", handlers.GetLinkByIDHandler)
	return r
}
func main() {
	log.Fatal(http.ListenAndServe(":8080", shortenerRouter()))
}
