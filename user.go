package goezyrouting

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func NewUserHandler(h *Handler) *UserHandler {
	uh := &UserHandler{h}
	return uh
}

type UserHandler struct{ *Handler }

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	// handle /user
	case "":
		switch r.Method {
		case "GET":
			h.usersGet()(w, r)
		case "POST":
			h.usersPost()(w, r)
		default:
			h.r.ErrMethodNotAllowed()(w, r)
		}
	// handle != "/"
	// like /user/:id or /user/:id/profile
	default:
		//if !h.ac.HasAccess(0, RoleUser){
		//	h.r.ErrUnauthorized(w,r)
		//	return
		//}
		// get id of param and save to context
		id, err := strconv.Atoi(head)
		if err != nil {
			h.r.Error(w, r, errors.New(fmt.Sprintf("Invalid user id %q", head)), usersViews)
			return
		}
		// do a db check here

		// save it to the context
		ctx = WithValue(r, ctx, "user_id", id)

		head, r.URL.Path = ShiftPath(r.URL.Path)
		switch head {
		// handle /user/:id
		case "":
			switch r.Method {
			case "GET":
				h.userGet(id)(w, r)
			case "PUT":
				h.handlePut(id)(w, r)
			case "DELETE":
				h.usersPost()(w, r)
			default:
				h.r.ErrMethodNotAllowed()(w, r)
			}
			// /user/:id/profile
		case "profile":
			ph := NewProfileHandler(h.Handler)
			h.CloseRoute(ph.userProfileGet()).ServeHTTP(w, r)
		// /user/:id/password
		case "password":
			// do something with the password
		default:
			h.r.ErrNotFound()(w, r)
		}
	}
}

const (
	usersViews = "user/index.html"
	userView   = "user/single.html"
)

func (h *UserHandler) usersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("Mission complete. Fetching Users...")
		uu := []User{
			{
				Name: "James Smith",
			},
			{
				Name: "Jane Smith",
			},
		}
		h.r.Render(w, r, http.StatusOK, uu, usersViews)
	}
}
func (h *UserHandler) usersPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("Mission complete. Adding Users...")
		h.r.Render(w, r, http.StatusCreated, map[string]interface{}{"success": "new user has been created"})
	}
}

type User struct {
	Name string `json:"name,omitempty"`
}

func (h *UserHandler) userGet(id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("Mission complete. Got ID:", id)
		u := new(User)
		u.Name = "James smith"
		ctx := r.Context()
		h.r.Render(w, r.WithContext(ctx), http.StatusOK, u, userView)
	}
}
func (h *UserHandler) handlePut(id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Log("user put handler", "ID:", id)
	}
}

func NewProfileHandler(h *Handler) *ProfileHandler {
	return &ProfileHandler{h}
}

type ProfileHandler struct{ *Handler }

func (h *ProfileHandler) userProfileGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in profile handler", r.Context().Value("user_id"))
		id := r.Context().Value("user_id")
		w.Write([]byte(fmt.Sprintf("profile for user %d", id)))
	}
}
