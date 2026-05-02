//go:build dev

package main

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.json
var openapiJSON string

func init() {
	openAPIHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(openapiJSON))
	}
}
