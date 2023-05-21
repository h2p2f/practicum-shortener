package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-shortener/internal/app/config"
	"github.com/h2p2f/practicum-shortener/internal/app/handler"
	"github.com/h2p2f/practicum-shortener/internal/app/logger"
	"github.com/h2p2f/practicum-shortener/internal/app/storage"
	"log"
	"net/http"
	"os"
	"strings"
)

//start up parameters
var runAddr, resultAddr string

//shortenerRouter creates a http router for two handlers
func shortenerRouter(s, r string) chi.Router {
	//create a storage and config
	stor := storage.NewLinkStorage()
	conf := config.NewServerConfig()
	conf.SetConfig(s, r)
	//message for app user
	message := fmt.Sprintf("Running Shortener. Server address: %s, Base URL: %s", s, r)
	fmt.Println(message)
	//create a router and add handlers
	handlers := handler.NewStorageHandler(stor, conf)
	c := chi.NewRouter()
	loggedRouter := c.With(logger.WithLogging)
	loggedRouter.Post("/", handlers.PostLinkHandler)
	loggedRouter.Get("/{id}", handlers.GetLinkByIDHandler)
	loggedRouter.Post("/api/shorten", handlers.PostLinkAPIHandler)
	return c
}
func main() {
	//get parameters from command line or environment variables
	flag.StringVar(&runAddr, "a", "localhost:8080", "address to run server on")
	flag.StringVar(&resultAddr, "b", "localhost:8080", "link to return")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envResultAddr := os.Getenv("BASE_URL"); envResultAddr != "" {
		resultAddr = envResultAddr
	}
	//cut protocol from resultAddr
	sliceAddr := strings.Split(resultAddr, "//")
	resultAddr = sliceAddr[len(sliceAddr)-1]
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	//start server
	log.Fatal(http.ListenAndServe(runAddr, shortenerRouter(runAddr, resultAddr)))
}
