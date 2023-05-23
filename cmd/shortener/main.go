package main

import (
	"context"
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
	"time"
)

// start up parameters
var runAddr, resultAddr, filePath string

// shortenerRouter creates a http router for two handlers
func shortenerRouter(stor *storage.LinkStorage, conf *config.ServerConfig, file *storage.FileDB) chi.Router {
	//create a storage and config
	//stor := storage.NewLinkStorage()
	//conf := config.NewServerConfig()
	//conf.SetConfig(s, r)
	//message for app user
	//message := fmt.Sprintf("Running Shortener. Server address: %s, Base URL: %s", s, r)
	//fmt.Println(message)
	//create a router and add handlers
	handlers := handler.NewStorageHandler(stor, conf, file)
	c := chi.NewRouter()
	//loggedRouter := c.With(logger.WithLogging)
	loggedAndZippedRouter := c.With(logger.WithLogging, handler.GzipHanle)
	loggedAndZippedRouter.Post("/", handlers.PostLinkHandler)
	loggedAndZippedRouter.Get("/{id}", handlers.GetLinkByIDHandler)
	loggedAndZippedRouter.Post("/api/shorten", handlers.PostLinkAPIHandler)
	return c
}
func main() {
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}

	//get parameters from command line or environment variables
	flag.StringVar(&runAddr, "a", "localhost:8080", "address to run server on")
	flag.StringVar(&resultAddr, "b", "localhost:8080", "link to return")
	flag.StringVar(&filePath, "f", "/tmp/short-url-db.json", "path to file with links")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envResultAddr := os.Getenv("BASE_URL"); envResultAddr != "" {
		resultAddr = envResultAddr
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		filePath = envFilePath
	}
	//cut protocol from resultAddr
	sliceAddr := strings.Split(resultAddr, "//")
	resultAddr = sliceAddr[len(sliceAddr)-1]
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	logger.Log.Sugar().Infof("Running Shortener. Server address: %s, Base URL: %s", runAddr, resultAddr)
	logger.Log.Sugar().Infof("File path: %s", filePath)
	stor := storage.NewLinkStorage()
	conf := config.NewServerConfig()
	conf.SetConfig(runAddr, resultAddr)

	storeInterval := 30 * time.Second
	fileDB := storage.NewFileDB(filePath, storeInterval, logger.Log)

	//read from file
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	data, err := fileDB.Read(ctx)
	if err != nil {
		fmt.Printf("error reading from file: %v", err)
	} else {
		stor.LoadAll(data)
	}

	//data = stor.GetAllSliced()
	//err = fileDB.Write(ctx, data)
	//if err != nil {
	//	fmt.Printf("error writing to file: %v", err)
	//}

	//go func() {
	//	for {
	//		time.Sleep(storeInterval)
	//		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//		defer cancel()
	//
	//		data := stor.GetAllSliced()
	//		//fmt.Println(data)
	//		err := fileDB.Write(ctx, data)
	//		if err != nil {
	//			fmt.Printf("error writing to file: %v", err)
	//		}
	//	}
	//}()

	//start server
	log.Fatal(http.ListenAndServe(runAddr, shortenerRouter(stor, conf, fileDB)))
}
