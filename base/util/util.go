package util

import (
	"encoding/json"
	"net/http"
)

func WritePayload(w http.ResponseWriter, status int, payload interface{}) {

	response, errorMarshal := json.MarshalIndent(payload, "", "  ")

	if errorMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMarshal.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func WriteError(w http.ResponseWriter, code int, message string) {

	WritePayload(w, code, map[string]interface{}{
		"code":  code,
		"error": message,
	})
}
