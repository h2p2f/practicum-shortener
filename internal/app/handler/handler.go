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
	SetConfig(s, r string)
	GetConfig() (string, string)
	GetResultAddress() string
}

type Filer interface {
	Read(ctx context.Context) ([][]byte, error)
	Write(ctx context.Context, links [][]byte) error
}

// handler for storage with config
type StorageHandler struct {
	storage Storager
	config  Configer
	file    Filer
}

// constructor of handler
func NewStorageHandler(storage Storager, config Configer, file Filer) *StorageHandler {
	return &StorageHandler{
		storage: storage,
		config:  config,
		file:    file,
	}
}

type originLink struct {
	Link string `json:"url"`
}

type shortLink struct {
	Link string `json:"result"`
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
			//break loop
			break
		}
	}
	s.SaveToDB()
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
			break
		}
	}
	s.SaveToDB()
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
