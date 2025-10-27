package bdd

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	ih "github.com/leandronowras/device-api/internal/http"
)

type apiWorld struct {
	server *httptest.Server
	resp   *http.Response
	body   []byte
	lastID string
}

func (w *apiWorld) theAPIIsRunning() error {
	r := chi.NewRouter()

	h := ih.NewHandler()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/devices", h.CreateDevice)
		r.Get("/devices", h.ListDevices)
		r.Get("/devices/{id}", h.GetDevice)
		r.Patch("/devices/{id}", h.UpdateDevice)
		r.Delete("/devices/{id}", h.DeleteDevice)
	})

	w.server = httptest.NewServer(r)
	return nil
}

func (w *apiWorld) stopServer() {
	if w.server != nil {
		w.server.Close()
	}
}
