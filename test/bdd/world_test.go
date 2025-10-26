package bdd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type apiWorld struct {
	server *httptest.Server
	resp   *http.Response
	body   []byte
	lastID string
}

func (w *apiWorld) theAPIIsRunning() error {
	// Minimal router for testing
	r := chi.NewRouter()

	r.Post("/devices", func(wr http.ResponseWriter, r *http.Request) {
		var in struct {
			Name  string `json:"name"`
			Brand string `json:"brand"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)

		out := map[string]any{
			"id":            "123",
			"name":          in.Name,
			"brand":         in.Brand,
			"state":         "available",
			"creation_time": "2025-10-26T00:00:00Z",
		}

		wr.Header().Set("Content-Type", "application/json")
		wr.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(wr).Encode(out)
	})

	w.server = httptest.NewServer(r)
	return nil
}

func (w *apiWorld) stopServer() {
	if w.server != nil {
		w.server.Close()
	}
}
