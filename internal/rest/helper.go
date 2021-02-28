package rest

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type Context struct {
	http.ResponseWriter
	status int
}

func (c *Context) GetStatus() int {
	return c.status
}

func (c *Context) WriteHeader(code int) {
	if c.status == 0 {
		c.ResponseWriter.WriteHeader(code)
		c.status = code
	}
}

func (c *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := c.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijack dont work")
	}

	return h.Hijack()
}

var v = validator.New()

func RespondJSON(w http.ResponseWriter, status int, obj interface{}) {
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

func RespondError(w http.ResponseWriter, status int, err error) {
	respondJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func BindJSON(r *http.Request, obj interface{}) error {
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
