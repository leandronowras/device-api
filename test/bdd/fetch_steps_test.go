package bdd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

// Given a device exists with name "iPhone" and brand "Apple"
func (w *apiWorld) aDeviceExistsWithNameAndBrand(name, brand string) error {
	payload := fmt.Sprintf(`{ "name": %q, "brand": %q }`, name, brand)
	if err := w.iPOSTWithJSON("/devices", &godog.DocString{Content: payload}); err != nil {
		return err
	}
	// expect creation to succeed so we can extract id
	if w.resp == nil || w.resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected 201 creating device, got %v (body=%s)", statusCode(w.resp), string(w.body))
	}
	// extract id from response body
	var m map[string]any
	if err := json.Unmarshal(w.body, &m); err != nil {
		return fmt.Errorf("invalid create json: %w", err)
	}
	id, _ := m["id"].(string)
	if id == "" {
		return fmt.Errorf("missing id in create response")
	}
	w.lastID = id
	return nil
}

// When I GET "/devices/{id}"
func (w *apiWorld) iGET(path string) error {
	if strings.Contains(path, "{id}") {
		if w.lastID == "" {
			return fmt.Errorf("no stored id to substitute into %q", path)
		}
		path = strings.ReplaceAll(path, "{id}", w.lastID)
	}
	resp, err := http.Get(w.server.URL + path)
	if err != nil {
		return err
	}
	w.resp = resp
	w.body, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return nil
}

// And the response json has keys: "id", "state", "creation_time"
func (w *apiWorld) theResponseJsonHasKeys(k1, k2, k3 string) error {
	var m map[string]any
	if err := json.Unmarshal(w.body, &m); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	for _, k := range []string{k1, k2, k3} {
		if _, ok := m[k]; !ok {
			return fmt.Errorf("missing key %q", k)
		}
	}
	return nil
}

// Given the API is running reacheable via http
func (w *apiWorld) theAPIIsRunningReacheableViaHttp() error {
	if w == nil || w.server == nil {
		return fmt.Errorf("test server not initialized")
	}
	return nil
}
