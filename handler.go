package goezyrouting

import (
	"net/http"
	"path"
	"strings"
)

func NewHandler(l Logger, r Renderer, ac AccessController) *Handler {
	return &Handler{
		l:  l,
		r:  r,
		ac: ac,
	}
}

type Handler struct {
	l  Logger
	r  Renderer
	ac AccessController
}

func Use(handler http.Handler, useMiddlewares ...middleware) http.Handler {
	// loop over the mws in reverse order
	// makes it more readable – middleware is called in the order it was added
	// use(h, mw1,mw2,mw3,...)
	for i := len(useMiddlewares) - 1; i >= 0; i-- {
		handler = useMiddlewares[i](handler)
	}
	return handler
}

// CloseRoute will close the route not allowing any additional paths
// example:
// Calling CloseRoute(someHandler) that is executed on path /user/1/profile
// will reject and respond "not found" if the requested path is /user/1/profile/1, /user/1/profile/foo, etc.
func (h *Handler) CloseRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if head, _ := ShiftPath(r.URL.Path); head != "" {
			h.r.ErrNotFound().ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}

	return p[1:i], p[i:]
}
