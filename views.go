package goezyrouting

import (
	"github.com/randallmlough/goezyrouting/website"
	"net/http"
)

const (
	forbiddenView   = "forbidden.html"
	badRequestView  = "400.html"
	notFoundView    = "404.html"
	serverErrorView = "500.html"
)

func NewViews(c *WebsiteConfig) (*View, error) {
	v, err := website.NewViews(
		c.TemplateDir,
		c.BaseTemplate,
		c.ViewsDir,
		c.Title,
		c.Styles,
	)
	if err != nil {
		return nil, err
	}
	return &View{v}, nil
}

type View struct {
	*website.View
}

func (v *View) RenderHTML(w http.ResponseWriter, r *http.Request, template string, data interface{}) {
	v.setOriginalPath(r)
	if err := v.View.Render(w, r, template, data); err != nil {
		v.ErrInternal()(w, r)
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, status int, data interface{}, template ...string) {
	v.setOriginalPath(r)
	if len(template) > 0 && template[0] != "" {
		if err := v.View.Render(w, r, template[0], data); err != nil {
			v.ErrInternal()(w, r)
		}
	} else {
		v.ErrInternal().ServeHTTP(w, r)
	}
}

func (v *View) Error(w http.ResponseWriter, r *http.Request, err error, template ...string) {
	d := &website.Page{
		Error: err,
	}
	v.setOriginalPath(r)
	if len(template) > 0 && template[0] != "" {
		if err := v.View.Render(w, r, template[0], d); err != nil {
			v.ErrInternal()(w, r)
		}
	} else {
		v.ErrInternal().ServeHTTP(w, r)
	}
}

func (v *View) ErrUnauthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.setHeader(w, http.StatusForbidden)
		d := v.error(renderErrForbidden)
		v.setOriginalPath(r)
		if err := v.View.Render(w, r, forbiddenView, d); err != nil {
			v.ErrInternal()(w, r)
		}
	}
}

func (v *View) ErrNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.setHeader(w, http.StatusNotFound)
		d := v.error(renderErrNotFound)
		v.setOriginalPath(r)
		if err := v.View.Render(w, r, notFoundView, d); err != nil {
			v.ErrInternal()(w, r)
		}
	}
}
func (v *View) ErrMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.ErrBadRequest()(w, r)
	}
}
func (v *View) ErrBadRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := v.error(renderErrBadRequest)
		v.setHeader(w, http.StatusBadRequest)
		v.setOriginalPath(r)
		if err := v.View.Render(w, r, badRequestView, d); err != nil {
			v.ErrInternal()(w, r)
		}
	}
}

func (v *View) ErrInternal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := v.error(renderErrGeneric)
		v.setHeader(w, http.StatusInternalServerError)
		v.setOriginalPath(r)
		v.View.Render(w, r, serverErrorView, d)
	}
}

func (v *View) setHeader(w http.ResponseWriter, status int) {
	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	w.WriteHeader(status)
}

func (v *View) setOriginalPath(r *http.Request) {
	ctx := r.Context()
	path := getURLPath(ctx)
	r.URL.Path = path
	r = r.WithContext(ctx)
}
func (v *View) error(err string) *website.Page {
	return &website.Page{
		Error: viewError(err),
	}
}

type viewError string

func (e viewError) Error() string {
	return string(e)
}
