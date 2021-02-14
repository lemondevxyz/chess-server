package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

var v = validator.New()

func respondJSON(w http.ResponseWriter, status int, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func respondError(w http.ResponseWriter, status int, err error) {
	respondJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func bindJSON(r *http.Request, obj interface{}) error {
	if r == nil || obj == nil || r.Body == nil {
		return errors.New("invalid parameters")
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(obj)
	if err != nil {
		return fmt.Errorf("json: %w", err)
	}

	err = v.Struct(obj)
	if err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	return nil
}
