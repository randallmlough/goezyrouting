package goezyrouting

import (
	"bytes"
	"encoding/gob"
	"net/http"
)

const (
	renderErrGeneric          = "something went wrong. Try again."
	renderErrMethodNotAllowed = "method not allowed"
	renderErrForbidden        = "forbidden"
	renderErrNotFound         = "not found"
	renderErrBadRequest       = "bad request"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, status int, data interface{}, template ...string)
	Error(w http.ResponseWriter, r *http.Request, err error, template ...string)
	ErrMethodNotAllowed() http.HandlerFunc
	ErrUnauthorized() http.HandlerFunc
	ErrNotFound() http.HandlerFunc
	ErrBadRequest() http.HandlerFunc
	ErrInternal() http.HandlerFunc
}

func NewRenderer() *Render {
	return &Render{}
}

type Render struct{}

func (*Render) Render(w http.ResponseWriter, r *http.Request, status int, data interface{}, template ...string) {
	b, err := ToBytes(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(b)
}

func ToBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (*Render) Error(w http.ResponseWriter, r *http.Request, err error, template ...string) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (*Render) ErrMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, renderErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (*Render) ErrUnauthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, renderErrForbidden, http.StatusUnauthorized)
	}
}
func (*Render) ErrNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, renderErrNotFound, http.StatusNotFound)
	}
}

func (*Render) ErrBadRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, renderErrBadRequest, http.StatusBadRequest)
	}
}

func (*Render) ErrInternal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, renderErrGeneric, http.StatusInternalServerError)
	}
}
