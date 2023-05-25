package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// string of letters for random string generation
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// random string generator
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// for future
//
//	func getIdFromBody(s string) string {
//		p := strings.Split(s, "/")
//		return p[len(p)-1]
//	}
//
// interface for storage
type Storager interface {
	Get(id string) (string, bool)
	Set(id, link string)
	Delete(id string)
	List() map[string]string
	Count() int
	GetAllSliced() [][]byte
}

// interface for config
type Configer interface {
	SetConfig(s, r string, f, d bool)
	GetConfig() (string, string, bool, bool)
	GetResultAddress() string
	UseDB() bool
	UseFile() bool
}

type Filer interface {
	Read(ctx context.Context) ([][]byte, error)
	Write(ctx context.Context, links [][]byte) error
}

type Databaser interface {
	PingContext(ctx context.Context) error
	InsertMetric(ctx context.Context, id string, oLink string) (err error)
}

// handler for storage with config
type StorageHandler struct {
	storage  Storager
	config   Configer
	file     Filer
	dataBase Databaser
}

// constructor of handler
func NewStorageHandler(storage Storager, config Configer, file Filer, db Databaser) *StorageHandler {
	return &StorageHandler{
		storage:  storage,
		config:   config,
		file:     file,
		dataBase: db,
	}
}

type originLink struct {
	Link string `json:"url"`
}

type shortLink struct {
	Link string `json:"result"`
}

type batchRequestLinks struct {
	CorrelationID string `json:"correlation_id"`
	OriginLink    string `json:"original_url"`
}

type batchResponseLinks struct {
	CorrelationID string `json:"correlation_id"`
	ShortLink     string `json:"short_url"`
}

// handler for get short link by request
func (s *StorageHandler) PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	//check method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//get config, set up mask for short link
	shortLink := "http://" + s.config.GetResultAddress() + "/"
	//read body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//check body
	if len(requestBody) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uuid := ""
	//generate random string and check if it is unique
	for {
		id := RandStringBytes(8)
		if _, ok := s.storage.Get(id); ok {
			//if not unique, generate new
			continue
		} else {
			//if unique, write to storage and return short link
			s.storage.Set(id, string(requestBody))
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(shortLink + id))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			uuid = id
			//break loop
			break
		}
	}
	if s.config.UseDB() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := s.dataBase.InsertMetric(ctx, uuid, string(requestBody))
		if err != nil {
			fmt.Println(err)
		}
	} else if s.config.UseFile() {
		s.SaveToDB()
	}
}

// GetLinkByIDHandler handler for get link by id
func (s *StorageHandler) GetLinkByIDHandler(w http.ResponseWriter, r *http.Request) {
	//check method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//get id from url
	id := chi.URLParam(r, "id")
	//check if id is in storage
	if link, ok := s.storage.Get(id); ok {
		//if yes, return 307 and redirect
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		//if not, return 404
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *StorageHandler) PostLinkAPIHandler(w http.ResponseWriter, r *http.Request) {
	//check method

	var buf bytes.Buffer
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sLink := "http://" + s.config.GetResultAddress() + "/"
	//read body
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//check body
	if buf.Len() == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	origin := originLink{Link: ""}
	err = json.Unmarshal(buf.Bytes(), &origin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	result := shortLink{""}
	uuid := ""
	//generate random string and check if it is unique
	for {
		id := RandStringBytes(8)
		if _, ok := s.storage.Get(id); ok {
			//if not unique, generate new
			continue
		} else {
			//if unique, write to storage and return short link
			s.storage.Set(id, origin.Link)
			result.Link = sLink + id
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			res, err := json.Marshal(result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = w.Write(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//break loop
			uuid = id
			break
		}
	}
	if s.config.UseDB() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := s.dataBase.InsertMetric(ctx, uuid, string(origin.Link))
		if err != nil {
			fmt.Println(err)
		}
	} else if s.config.UseFile() {
		s.SaveToDB()
	}
}

func (s *StorageHandler) SaveToDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	data := s.storage.GetAllSliced()
	//fmt.Println(data)
	err := s.file.Write(ctx, data)
	if err != nil {
		fmt.Printf("error writing to file: %v", err)
	}
}

func (s *StorageHandler) DBPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.dataBase.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return
	}
}

func (s *StorageHandler) BatchMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var buf bytes.Buffer
	var reqLink []batchRequestLinks
	var resLinks []batchResponseLinks

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if buf.Len() == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &reqLink)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	sLink := "http://" + s.config.GetResultAddress() + "/"
	for _, link := range reqLink {
		uuid := ""
		//generate random string and check if it is unique
		for {
			id := RandStringBytes(8)
			if _, ok := s.storage.Get(id); ok {
				//if not unique, generate new
				continue
			} else {
				//if unique, write to storage and return short link
				s.storage.Set(id, link.OriginLink)
				resLinks = append(resLinks, batchResponseLinks{
					link.CorrelationID,
					sLink + id,
				})

			}
			//break loop
			uuid = id
			break
		}
		if s.config.UseDB() {
			err := s.dataBase.InsertMetric(ctx, uuid, link.OriginLink)
			if err != nil {
				fmt.Println(err)
			}
		} else if s.config.UseFile() {
			s.storage.Set(uuid, link.OriginLink)
			s.SaveToDB()
		}
	}

	data, err := json.Marshal(resLinks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}
