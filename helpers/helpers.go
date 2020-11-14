package helpers

import (
	"encoding/json"
	"net/http"
)

// ServerError Model ...
type serverError struct {
	Message string `json:"message"`
}

// RespondWithError ...
func RespondWithError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorMessage)
}

// ResponseJSON ...
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
