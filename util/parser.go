package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// parseRequestBody parses request body to given interface
func ParseRequestBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(v)
	default:
		return fmt.Errorf("error: Unsupported Content-Type: %s", contentType)
	}
}
