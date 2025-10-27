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
	if err := w.iPOSTWithJSON("/v1/devices", &godog.DocString{Content: payload}); err != nil {
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

// Given there are more than {n} devices stored
func (w *apiWorld) thereAreMoreThanDevicesStored(n int) error {
	target := n + 1
	for i := 0; i < target; i++ {
		payload := fmt.Sprintf(`{ "name": "dev-%d", "brand": "brand-%d" }`, i, i)
		if err := w.iPOSTWithJSON("/v1/devices", &godog.DocString{Content: payload}); err != nil {
			return fmt.Errorf("seed POST failed at %d: %w", i, err)
		}
		if w.resp == nil || w.resp.StatusCode != 201 {
			return fmt.Errorf("seed create expected 201, got %d (body=%s)", statusCode(w.resp), string(w.body))
		}
	}
	return nil
}

// Then the response json should contain {n} devices
func (w *apiWorld) theResponseJSONShouldContainNDevices(n int) error {
	var anyJSON any
	if err := json.Unmarshal(w.body, &anyJSON); err != nil {
		return fmt.Errorf("invalid json: %w; body=%s", err, string(w.body))
	}

	switch v := anyJSON.(type) {
	case []any:
		if len(v) != n {
			return fmt.Errorf("expected %d devices in array, got %d", n, len(v))
		}
		return nil
	case map[string]any:
		// common shape: { "items": [...], "next_page": "...", "previous_page": "..." }
		items, ok := v["items"].([]any)
		if !ok {
			// also allow "data" as a fallback
			items, ok = v["data"].([]any)
			if !ok {
				return fmt.Errorf(`json does not contain an "items" (or "data") array`)
			}
		}
		if len(items) != n {
			return fmt.Errorf("expected %d devices in items, got %d", n, len(items))
		}
		return nil
	default:
		return fmt.Errorf("unexpected json root type %T", v)
	}
}

// And the response json should include "next_page" and "previous_page" fields
func (w *apiWorld) theResponseJSONShouldIncludeNextPrev() error {
	var obj map[string]any
	if err := json.Unmarshal(w.body, &obj); err != nil {
		return fmt.Errorf("invalid json: %w; body=%s", err, string(w.body))
	}
	if _, ok := obj["next_page"]; !ok {
		return fmt.Errorf(`missing "next_page" field`)
	}
	if _, ok := obj["previous_page"]; !ok {
		return fmt.Errorf(`missing "previous_page" field`)
	}
	return nil
}

func theAPIIsRunning() error {
	return nil
}

// Then the response json at "{jsonpath}" should be "{expected}"
func (w *apiWorld) responseJsonAtShouldBe(path, expected string) error {
	var body any
	if err := json.Unmarshal(w.body, &body); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Support "$[0].field" and "$.field"
	if strings.HasPrefix(path, "$[0].") {
		field := strings.TrimPrefix(path, "$[0].")
		arr, ok := body.([]any)
		if !ok || len(arr) == 0 {
			return fmt.Errorf("expected array at root for %s", path)
		}
		obj, ok := arr[0].(map[string]any)
		if !ok {
			return fmt.Errorf("array element 0 is not an object")
		}
		val := fmt.Sprintf("%v", obj[field])
		if val != expected {
			return fmt.Errorf("expected %s=%q, got %q", field, expected, val)
		}
		return nil
	}

	if strings.HasPrefix(path, "$.") {
		field := strings.TrimPrefix(path, "$.")
		obj, ok := body.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object at root for %s", path)
		}
		val := fmt.Sprintf("%v", obj[field])
		if val != expected {
			return fmt.Errorf("expected %s=%q, got %q", field, expected, val)
		}
		return nil
	}

	return fmt.Errorf("unsupported JSON path: %s", path)
}
