package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-shortener/internal/app/config"
	"github.com/h2p2f/practicum-shortener/internal/app/storage"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//for future
//func getIdFromBody(s string) string {
//	p := strings.Split(s, "/")
//	return p[len(p)-1]
//}

type StorageHandler struct {
	storage storage.Storage
	config  config.Config
}

func NewStorageHandler(storage storage.Storage, config config.Config) *StorageHandler {
	return &StorageHandler{
		storage: storage,
		config:  config,
	}
}
func (s *StorageHandler) PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	shortLink := "http://" + s.config.GetResultAddress() + "/"
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(requestBody) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := RandStringBytes(8)
	if _, ok := s.storage.Get(id); ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		s.storage.Set(id, string(requestBody))
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(shortLink + id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *StorageHandler) GetLinkByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParam(r, "id")

	if link, ok := s.storage.Get(id); ok {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
