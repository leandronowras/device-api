package bdd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cucumber/godog"
)

func (w *apiWorld) iPOSTWithJSON(path string, doc *godog.DocString) error {
	resp, err := http.Post(w.server.URL+path, "application/json", bytes.NewBufferString(doc.Content))
	if err != nil {
		return err
	}
	w.resp = resp
	w.body, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return nil
}

func (w *apiWorld) theResponseCodeShouldBe(code int) error {
	if w.resp == nil || w.resp.StatusCode != code {
		return fmt.Errorf("expected %d, got %v (body=%s)", code, statusCode(w.resp), string(w.body))
	}
	return nil
}

func (w *apiWorld) jsonAtShouldBe(path, expected string) error {
	key, err := topLevelKeyFromPath(path)
	if err != nil {
		return err
	}
	var m map[string]any
	if err := json.Unmarshal(w.body, &m); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	got, ok := m[key]
	if !ok {
		return fmt.Errorf("missing key %q", key)
	}
	if fmt.Sprint(got) != expected {
		return fmt.Errorf(`json at %s: want %q, got %q`, path, expected, fmt.Sprint(got))
	}
	return nil
}

func theResponseJsonHasKeys(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}
