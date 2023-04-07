package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-shortener/internal/app/config"
	"github.com/h2p2f/practicum-shortener/internal/app/handler"
	"github.com/h2p2f/practicum-shortener/internal/app/storage"
	"log"
	"net/http"
)

var runAddr string
var resultAddr string

func shortenerRouter(s, r string) chi.Router {
	stor := storage.NewLinkStorage()
	conf := config.NewServerConfig()
	conf.SetConfig(s, r)
	fmt.Println(conf)
	handlers := handler.NewStorageHandler(stor, conf)
	c := chi.NewRouter()
	c.Post("/", handlers.PostLinkHandler)
	c.Get("/{id}", handlers.GetLinkByIDHandler)
	return c
}
func main() {

	flag.StringVar(&runAddr, "a", "localhost:8080", "address to run server on")
	flag.StringVar(&resultAddr, "b", "localhost:8080", "link to return")
	flag.Parse()

	log.Fatal(http.ListenAndServe(runAddr, shortenerRouter(runAddr, resultAddr)))
}
