package goezyrouting

import "net/http"

func NewAuthHandler(h *Handler) *AuthHandler {
	return &AuthHandler{h}
}

type AuthHandler struct{ *Handler }

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "login":
		switch r.Method {
		case "POST":
			h.loginPost()(w, r)
		default:
			h.r.ErrMethodNotAllowed()(w, r)
		}
	case "register":
		switch r.Method {
		case "POST":
		default:
			h.r.ErrMethodNotAllowed()(w, r)
		}
	case "logout":
		switch r.Method {
		case "POST":
			requireUser(h.logoutPost()).ServeHTTP(w, r)
		default:
			h.r.ErrMethodNotAllowed()(w, r)
		}
	default:
		h.r.ErrNotFound()(w, r)
	}
}

func (h *AuthHandler) loginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("mission complete.", "responding to request")
		w.Write([]byte("you are logged in"))
	}
}

func (h *AuthHandler) logoutPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("mission complete.", "responding to request")
		w.Write([]byte("you are logged out"))
	}
}
