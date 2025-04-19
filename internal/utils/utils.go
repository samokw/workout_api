package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]interface{}

func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ExtractParam(param string, w http.ResponseWriter, r *http.Request) (int64, error) {
	paramID := chi.URLParam(r, param)
	if paramID == "" {
		http.Error(w, fmt.Sprintf("Missing URL parameter: %s", param), http.StatusBadRequest)
		return 0, fmt.Errorf("missing URL parameter: %s", param)
	}
	paramIntID, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URL parameter format for %s: must be an integer", param), http.StatusBadRequest)
		return 0, fmt.Errorf("invalid URL parameter format for %s: %w", param, err)
	}
	return paramIntID, nil
}
