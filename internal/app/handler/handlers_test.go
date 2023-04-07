package handler

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-shortener/internal/app/config"
	"github.com/h2p2f/practicum-shortener/internal/app/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStorageHandler_PostLinkHandler(t *testing.T) {
	type want struct {
		statusCode int
		shortLink  string
	}
	tests := []struct {
		name   string
		method string
		link   string
		want   want
	}{
		{
			name:   " Positive test",
			method: http.MethodPost,
			link:   "https://google.com",
			want: want{
				statusCode: 201,
				shortLink:  "http://localhost:8080/12345678",
			},
		},
		{
			name:   " Negative test",
			method: http.MethodPost,
			link:   "",
			want: want{
				statusCode: 400,
				shortLink:  "",
			},
		},
		{
			name:   " Negative test 2",
			method: http.MethodGet,
			link:   "https://google.com",
			want: want{
				statusCode: 405,
				shortLink:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.link))
			r := chi.NewRouter()
			s := storage.NewLinkStorage()
			c := config.NewServerConfig()
			c.SetConfig("localhost:8080", "localhost:8080")
			handlers := NewStorageHandler(s, c)
			if tt.method == http.MethodPost {
				r.Post("/", handlers.PostLinkHandler)
			} else {
				r.Get("/", handlers.PostLinkHandler)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.want.statusCode {
				t.Errorf("PostLinkHandler() got = %v, want %v", w.Code, tt.want.statusCode)
			}
		})
	}
}

//rewrite this code!!!!!!!!!!!!!
func TestStorageHandler_GetLinkByIDHandler(t *testing.T) {
	type want struct {
		statusCode int
		link       string
	}
	tests := []struct {
		name   string
		method string
		link   string
		id     string
		want   want
	}{
		{
			name:   " Positive test",
			method: http.MethodGet,
			id:     "",
			link:   "https://google.com",
			want: want{
				statusCode: 404,
				link:       "http://localhost:8080",
			},
		},
		{
			name:   " Negative test",
			method: http.MethodGet,
			id:     "",
			link:   "",
			want: want{
				statusCode: 404,
				link:       "",
			},
		},
		{
			name:   " Negative test 2",
			method: http.MethodPost,
			link:   "https://google.com",
			want: want{
				statusCode: 404,
				link:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.link))
			r := chi.NewRouter()
			s := storage.NewLinkStorage()
			c := config.NewServerConfig()
			c.SetConfig("localhost:8080", "localhost:8080")
			handlers := NewStorageHandler(s, c)
			req.Header.Set("Content-Type", "text/plain")
			postReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.link)))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, postReq)
			//res := w.Result()
			//body, _ := io.ReadAll(res.Body)
			//id := strings.Split(string(body), "/")
			//r.Post("/", handlers.PostLinkHandler)
			if tt.method == http.MethodGet {
				r.Get("/{id}", handlers.GetLinkByIDHandler)
			} else {
				r.Post("/{id}", handlers.GetLinkByIDHandler)
			}
			aw := httptest.NewRecorder()
			r.ServeHTTP(aw, req)
			if w.Code != tt.want.statusCode {
				t.Errorf("GetLinkByIDHandler() got = %v, want %v", w.Code, tt.want.statusCode)
			}

		})
	}
}
