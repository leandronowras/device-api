package http

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leandronowras/device-api/internal/device"
)

type Handler struct {
	devices map[string]*device.Device
}

func NewHandler() *Handler {
	return &Handler{devices: make(map[string]*device.Device)}
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

	h.devices[d.ID()] = d
	writeJSON(w, stdhttp.StatusCreated, toResp(d))
}

// --- READ (GET by ID) --------------------------------------------------------

func (h *Handler) GetDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	d, ok := h.devices[id]
	if !ok {
		writeJSONError(w, &device.DomainError{
			Code: "not_found", Field: "id", Message: "device not found", HTTP: stdhttp.StatusNotFound,
		})
		return
	}
	writeJSON(w, stdhttp.StatusOK, toResp(d))
}

// --- LIST (GET all or filtered) ---------------------------------------------

func (h *Handler) ListDevices(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	brand := strings.TrimSpace(r.URL.Query().Get("brand"))
	state := strings.TrimSpace(r.URL.Query().Get("state"))

	resp := []deviceResponse{}
	for _, d := range h.devices {
		if brand != "" && !strings.EqualFold(brand, d.Brand()) {
			continue
		}
		if state != "" && !strings.EqualFold(state, d.State()) {
			continue
		}
		resp = append(resp, toResp(d))
	}
	writeJSON(w, stdhttp.StatusOK, resp)
}

// --- UPDATE (PATCH minimal example) -----------------------------------------

func (h *Handler) UpdateDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	d, ok := h.devices[id]
	if !ok {
		writeJSONError(w, &device.DomainError{
			Code: "not_found", Field: "id", Message: "device not found", HTTP: stdhttp.StatusNotFound,
		})
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
		newD, err := device.New(*req.Name, d.Brand(), d.State())
		if err != nil {
			writeJSONError(w, err)
			return
		}
		*d = *newD // overwrite
	}
	if req.Brand != nil && strings.TrimSpace(*req.Brand) != "" {
		newD, err := device.New(d.Name(), *req.Brand, d.State())
		if err != nil {
			writeJSONError(w, err)
			return
		}
		*d = *newD
	}
	if req.State != nil && strings.TrimSpace(*req.State) != "" {
		newD, err := device.New(d.Name(), d.Brand(), *req.State)
		if err != nil {
			writeJSONError(w, err)
			return
		}
		*d = *newD
	}

	h.devices[id] = d
	writeJSON(w, stdhttp.StatusOK, toResp(d))
}

// --- DELETE ------------------------------------------------------------------

func (h *Handler) DeleteDevice(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id := chi.URLParam(r, "id")
	d, ok := h.devices[id]
	if !ok {
		writeJSONError(w, &device.DomainError{
			Code: "not_found", Field: "id", Message: "device not found", HTTP: stdhttp.StatusNotFound,
		})
		return
	}

	if d.State() == device.StateInUse {
		writeJSONError(w, device.ErrConflict("device", "cannot delete device in use"))
		return
	}

	delete(h.devices, id)
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
