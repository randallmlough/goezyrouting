package goezyrouting

import (
	"fmt"
	"net/http"
)

const (
	homeView = "index.html"
)

func NewStaticHandler(l Logger, r Renderer) *StaticHandler {
	h := new(StaticHandler)
	h.l = l
	h.r = r
	return h
}

type StaticHandler Handler

func (h *StaticHandler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("rendering template:", homeView)

		h.r.Render(w, r, http.StatusFound, nil, homeView)
	}
}
func (h *StaticHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("I'm login")
		w.Write([]byte("I'm going to log you in"))

	}
}
func (h *StaticHandler) Handler(id int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in profile handler", r.Context().Value("user_id"))
	})
}
