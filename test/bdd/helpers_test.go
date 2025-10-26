package bdd

import (
	"fmt"
	"net/http"
	"strings"
)

func topLevelKeyFromPath(path string) (string, error) {
	if !strings.HasPrefix(path, "$.") {
		return "", fmt.Errorf("only top-level paths like '$.name' are supported, got %q", path)
	}
	key := strings.TrimPrefix(path, "$.")
	if key == "" || strings.Contains(key, ".") {
		return "", fmt.Errorf("nested paths not supported in this minimal helper: %q", path)
	}
	return key, nil
}

func statusCode(r *http.Response) any {
	if r == nil {
		return nil
	}
	return r.StatusCode
}
