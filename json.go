package goezyrouting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func NewJson() *Json {
	return &Json{}
}

type Json struct {
	data interface{}
}

const (
	jsonErrGeneric          = `{"error":"` + renderErrGeneric + `"}`
	jsonErrMethodNotAllowed = `{"error":"` + renderErrMethodNotAllowed + `"}`
	jsonErrForbidden        = `{"error":"` + renderErrForbidden + `"}`
	jsonErrNotFound         = `{"error":"` + renderErrNotFound + `"}`
	jsonErrBadRequest       = `{"error":"` + renderErrBadRequest + `"}`
)

func (j *Json) Render(w http.ResponseWriter, r *http.Request, status int, data interface{}, template ...string) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		j.ErrInternal().ServeHTTP(w, r)
		return
	}
	j.setHeader(w, status)
	j.write(w, b)
}

func (j *Json) Error(w http.ResponseWriter, r *http.Request, err error, template ...string) {
	b, err := j.error(err)
	if err != nil {
		b = []byte(jsonErrGeneric)
	}
	j.setHeader(w, http.StatusBadRequest)
	j.write(w, b)
	return
}
func (j *Json) ErrMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j.setHeader(w, http.StatusMethodNotAllowed)
		j.write(w, []byte(jsonErrMethodNotAllowed))
		return
	}
}
func (j *Json) ErrUnauthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j.setHeader(w, http.StatusUnauthorized)
		j.write(w, []byte(jsonErrForbidden))
		return
	}
}
func (j *Json) ErrNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j.setHeader(w, http.StatusNotFound)
		j.write(w, []byte(jsonErrNotFound))
		return
	}
}
func (j *Json) ErrBadRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j.setHeader(w, http.StatusBadRequest)
		j.write(w, []byte(jsonErrBadRequest))
		return
	}
}
func (j *Json) ErrInternal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j.setHeader(w, http.StatusInternalServerError)
		j.write(w, []byte(jsonErrGeneric))
		return
	}
}

// helper funcs
func (j *Json) MarshalJSON() ([]byte, error) {
	fmt.Println("in own marshaller")
	return json.Marshal(j.data)
}

func (j *Json) setHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func (j *Json) write(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		log.Println(err)
		return
	}
}

// error into serialized object
func (j *Json) error(err error) ([]byte, error) {
	var tmp = new(struct {
		Error string `json:"error"`
	})
	tmp.Error = err.Error()

	b, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}
	return b, nil
}
