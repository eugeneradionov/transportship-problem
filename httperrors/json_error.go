package httperrors

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrJSON struct {
	Error string `json:"error,omitempty"`
	Code  int    `json:"code,omitempty"`
}

func SendErrorJSON(w http.ResponseWriter, err error, code int) {
	log.Printf("[DEBUG] Error: %s; Code: %d", err.Error(), code)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ErrJSON{Error: err.Error(), Code: code})
}
