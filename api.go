package goezyrouting

import (
	"errors"
	"net/http"
)

const apiV1 = "v1"

func NewAPI(l Logger, r Renderer, ac AccessController) *API {
	api := &API{
		Handler: &Handler{
			l:  l,
			r:  r,
			ac: ac,
		},
	}
	return api
}

type API struct {
	*Handler
}

func (h *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case apiV1:
		v1 := &ApiVersion1{h.Handler}
		v1.ServeHTTP(w, r)
	default:
		h.r.Error(w, r, errors.New("unrecognized API version"))
	}
}

type ApiVersion1 struct {
	*Handler
}

func (h *ApiVersion1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "auth":
		auth := NewAuthHandler(h.Handler)
		auth.ServeHTTP(w, r)
	case "user":
		uh := Use(
			NewUserHandler(h.Handler),
			minAccessLevel(3),
			levelTwoMiddleware,
		)
		uh.ServeHTTP(w, r)
	case "account":
	default:
		h.r.ErrNotFound()(w, r)
	}

}
