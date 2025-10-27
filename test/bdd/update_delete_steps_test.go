package bdd

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

func (w *apiWorld) iPATCHWithJSON(path string, doc *godog.DocString) error {
	url := w.server.URL + strings.Replace(path, "{id}", w.lastID, -1)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBufferString(doc.Content))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	w.resp = resp
	w.body, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return nil
}

func (w *apiWorld) iDELETE(path string) error {
	url := w.server.URL + strings.Replace(path, "{id}", w.lastID, -1)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	w.resp = resp
	w.body, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return nil
}
