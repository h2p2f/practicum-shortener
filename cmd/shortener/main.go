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
var runAddr, resultAddr, filePath, databaseVar string
var useDB, useFile bool

func isFlagPassed(s string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == s {
			found = true
		}
	})
	return found
}

// shortenerRouter creates a http router for two handlers
func shortenerRouter(stor *storage.LinkStorage, conf *config.ServerConfig, file *storage.FileDB, db *storage.PGDB) chi.Router {
	//create a storage and config
	//stor := storage.NewLinkStorage()
	//conf := config.NewServerConfig()
	//conf.SetConfig(s, r)
	//message for app user
	//message := fmt.Sprintf("Running Shortener. Server address: %s, Base URL: %s", s, r)
	//fmt.Println(message)
	//create a router and add handlers
	handlers := handler.NewStorageHandler(stor, conf, file, db)
	c := chi.NewRouter()
	//loggedRouter := c.With(logger.WithLogging)
	loggedAndZippedRouter := c.With(logger.WithLogging, handler.GzipHanle)
	loggedAndZippedRouter.Post("/", handlers.PostLinkHandler)
	loggedAndZippedRouter.Get("/{id}", handlers.GetLinkByIDHandler)
	loggedAndZippedRouter.Post("/api/shorten", handlers.PostLinkAPIHandler)
	loggedAndZippedRouter.Get("/ping", handlers.DBPing)
	loggedAndZippedRouter.Post("/api/shorten/batch", handlers.BatchMetrics)
	return c
}
func main() {
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	useDB = false
	useFile = false
	//get parameters from command line or environment variables
	flag.StringVar(&runAddr, "a", "localhost:8080", "address to run server on")
	flag.StringVar(&resultAddr, "b", "localhost:8080", "link to return")
	flag.StringVar(&filePath, "f", "/tmp/short-url-db.json", "path to file with links")
	flag.StringVar(&databaseVar, "d",
		"postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
		"databaseVar to store metrics")
	flag.Parse()

	if isFlagPassed("d") {
		useDB = true
	}

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envResultAddr := os.Getenv("BASE_URL"); envResultAddr != "" {
		resultAddr = envResultAddr
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		filePath = envFilePath
		useFile = true
	}
	if envDatabaseVar := os.Getenv("DATABASE_DSN"); envDatabaseVar != "" {
		databaseVar = envDatabaseVar
		useDB = true
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
	conf.SetConfig(runAddr, resultAddr, useDB, useFile)

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

	db := storage.NewPostgresDB(databaseVar, logger.Log)

	err = db.Create(ctx)
	if err != nil {
		fmt.Printf("error creating table: %v", err)
	}
	logger.Log.Sugar().Infof("Database: %s", databaseVar)
	logger.Log.Sugar().Infof("Use database: %t", useDB)
	logger.Log.Sugar().Infof("Use file: %t", useFile)
	logger.Log.Sugar().Infof("Store interval: %s", storeInterval)
	//start server
	log.Fatal(http.ListenAndServe(runAddr, shortenerRouter(stor, conf, fileDB, db)))
}
