package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/marcboeker/go-duckdb"

	ih "github.com/leandronowras/device-api/internal/http"
	duckdbrepo "github.com/leandronowras/device-api/internal/repository/duckdb"
)

func main() {
	db, err := sql.Open("duckdb", "./devices.db")
	if err != nil {
		log.Fatalf("failed to open duckdb: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping duckdb: %v", err)
	}

	repo := duckdbrepo.NewDeviceRepository(db)
	h := ih.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/devices", h.CreateDevice)
		r.Get("/devices", h.ListDevices)
		r.Get("/devices/{id}", h.GetDevice)
		r.Patch("/devices/{id}", h.UpdateDevice)
		r.Delete("/devices/{id}", h.DeleteDevice)
	})

	addr := ":8080"
	log.Printf("🚀 Device API running at http://localhost%s/v1/devices", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
