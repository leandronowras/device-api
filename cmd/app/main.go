package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	ih "github.com/leandronowras/device-api/internal/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	// Create the in-memory handler (no repository)
	h := ih.NewHandler()

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
