package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	var buffer bytes.Buffer

	if err := json.NewEncoder(&buffer).Encode(data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write(buffer.Bytes())
	return err
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	response := map[string]string{
		"error": message,
	}

	return WriteJSON(w, status, response)
}
