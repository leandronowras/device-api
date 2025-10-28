package bdd

import (
	"database/sql"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	_ "github.com/marcboeker/go-duckdb"

	ih "github.com/leandronowras/device-api/internal/http"
	duckdbrepo "github.com/leandronowras/device-api/internal/repository/duckdb"
)

type apiWorld struct {
	server *httptest.Server
	resp   *http.Response
	body   []byte
	lastID string
	db     *sql.DB
}

func (w *apiWorld) theAPIIsRunning() error {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return err
	}
	w.db = db

	repo := duckdbrepo.NewDeviceRepository(db)

	r := chi.NewRouter()
	h := ih.NewHandler(repo)

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
	if w.db != nil {
		w.db.Close()
	}
}
