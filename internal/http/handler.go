package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leandronowras/device-api/internal/device"
	"github.com/leandronowras/device-api/internal/repository"
)

type Handler struct {
	repo repository.DeviceRepository
}

func NewHandler(repo repository.DeviceRepository) *Handler {
	return &Handler{repo: repo}
}

// Shared response struct
type deviceResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Brand        string    `json:"brand"`
	State        string    `json:"state"`
	CreationTime time.Time `json:"creation_time"`
}

// Helper to convert domain to response
func toResp(d *device.Device) deviceResponse {
	return deviceResponse{
		ID:           d.ID(),
		Name:         d.Name(),
		Brand:        d.Brand(),
		State:        d.State(),
		CreationTime: d.CreationTime(),
	}
}

// --- CREATE ------------------------------------------------------------------

func (h *Handler) CreateDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req struct {
		Name  string `json:"name"`
		Brand string `json:"brand"`
		State string `json:"state,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, &device.DomainError{
			Code: "invalid_json", Message: "invalid JSON body", HTTP: stdhttp.StatusBadRequest,
		})
		return
	}

	var (
		d   *device.Device
		err error
	)
	if strings.TrimSpace(req.State) == "" {
		d, err = device.New(req.Name, req.Brand)
	} else {
		d, err = device.New(req.Name, req.Brand, req.State)
	}
	if err != nil {
		writeJSONError(w, err)
		return
	}

	saved, err := h.repo.Save(context.Background(), d)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusCreated, toResp(saved))
}

// --- READ (GET by ID) --------------------------------------------------------

func (h *Handler) GetDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	d, err := h.repo.FindByID(context.Background(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, &device.DomainError{
			Code: "not_found", Field: "id", Message: "device not found", HTTP: stdhttp.StatusNotFound,
		})
		return
	}
	if err != nil {
		writeJSONError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, toResp(d))
}

// --- LIST (GET all or filtered) ---------------------------------------------

func (h *Handler) ListDevices(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	brand := strings.TrimSpace(r.URL.Query().Get("brand"))
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var bPtr, sPtr *string
	if brand != "" {
		bPtr = &brand
	}
	if state != "" {
		sPtr = &state
	}

	list, err := h.repo.FindAll(context.Background(), bPtr, sPtr)
	if err != nil {
		writeJSONError(w, err)
		return
	}

	resp := []deviceResponse{}
	for _, d := range list {
		resp = append(resp, toResp(d))
	}

	if pageStr != "" || limitStr != "" {
		page := 1
		if pageStr != "" {
			if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
				page = int(p)
			}
		}
		
		limit := 10
		if limitStr != "" {
			if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 {
				limit = int(l)
				if limit > 100 {
					limit = 100
				}
			}
		}

		total := len(resp)
		start := (page - 1) * limit
		end := start + limit

		if start >= total {
			start = total
		}
		if end > total {
			end = total
		}

		paged := resp[start:end]
		
		nextPage := ""
		if end < total {
			nextPage = strconv.FormatInt(int64(page+1), 10)
		}
		
		prevPage := ""
		if page > 1 {
			prevPage = strconv.FormatInt(int64(page-1), 10)
		}

		envelope := map[string]any{
			"items":         paged,
			"next_page":     nextPage,
			"previous_page": prevPage,
		}
		writeJSON(w, stdhttp.StatusOK, envelope)
		return
	}

	writeJSON(w, stdhttp.StatusOK, resp)
}

// --- UPDATE (PATCH minimal example) -----------------------------------------

func (h *Handler) UpdateDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	d, err := h.repo.FindByID(context.Background(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, &device.DomainError{
			Code: "not_found", Field: "id", Message: "device not found", HTTP: stdhttp.StatusNotFound,
		})
		return
	}
	if err != nil {
		writeJSONError(w, err)
		return
	}

	var req struct {
		Name  *string `json:"name,omitempty"`
		Brand *string `json:"brand,omitempty"`
		State *string `json:"state,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, &device.DomainError{
			Code: "invalid_json", Message: "invalid JSON body", HTTP: stdhttp.StatusBadRequest,
		})
		return
	}

	// Business rule: cannot change name/brand if in-use
	if d.State() == device.StateInUse && (req.Name != nil || req.Brand != nil) {
		writeJSONError(w, device.ErrForbiddenChange("name/brand", "device is in use", stdhttp.StatusBadRequest))
		return
	}

	// Apply updates if provided
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		if err := d.SetName(*req.Name); err != nil {
			writeJSONError(w, err)
			return
		}
	}
	if req.Brand != nil && strings.TrimSpace(*req.Brand) != "" {
		if err := d.SetBrand(*req.Brand); err != nil {
			writeJSONError(w, err)
			return
		}
	}
	if req.State != nil && strings.TrimSpace(*req.State) != "" {
		if err := d.SetState(*req.State); err != nil {
			writeJSONError(w, err)
			return
		}
	}

	updated, err := h.repo.Update(context.Background(), d)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, toResp(updated))
}

// --- DELETE ------------------------------------------------------------------

func (h *Handler) DeleteDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	d, err := h.repo.FindByID(context.Background(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, &device.DomainError{
			Code: "not_found", Field: "id", Message: "device not found", HTTP: stdhttp.StatusNotFound,
		})
		return
	}
	if err != nil {
		writeJSONError(w, err)
		return
	}

	if d.State() == device.StateInUse {
		writeJSONError(w, device.ErrConflict("device", "cannot delete device in use"))
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		writeJSONError(w, err)
		return
	}
	w.WriteHeader(stdhttp.StatusNoContent)
}

// --- Helpers -----------------------------------------------------------------

func writeJSON(w stdhttp.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w stdhttp.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var derr *device.DomainError
	if errors.As(err, &derr) {
		w.WriteHeader(derr.HTTP)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    derr.Code,
			"field":   derr.Field,
			"message": derr.Message,
		})
		return
	}

	w.WriteHeader(stdhttp.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    "internal_error",
		"message": "unexpected error",
	})
}
