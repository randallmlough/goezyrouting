package goezyrouting

import (
	"net/http"
)

func Initialize(cfgFile string) (*Server, error) {
	c, err := LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	v, err := NewViews(c.Website)
	if err != nil {
		return nil, err
	}

	a := NewApp(
		WithConfig(c),
	)

	l := NewLogger()

	s := NewServer(
		DefaultServer(),
		WithErrorLog(l.ErrorLog()),
		WithPort(c.Application.Port),
		WithLogger(l),
		WithHandler(
			Use(
				a.AppRouter(l, NewRenderer(), v),
				Request,
				RequestID,
				Recoverer(l),
				logger(l),
				Cors,
				SecureHeaders,
			),
		),
	)

	return s, nil
}

func NewApp(opts ...Options) *App {
	a := new(App)
	for _, opt := range opts {
		opt(a)
	}
	return a
}

type App struct {
	config *Config
}

type Options func(*App)

func WithConfig(c *Config) Options {
	return func(a *App) {
		a.config = c
	}
}

func (a *App) AppRouter(l Logger, ren Renderer, v *View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		head, tail := ShiftPath(r.URL.Path)
		switch head {
		case "api":
			r.URL.Path = tail
			api := Use(
				NewAPI(l, NewJson(), NewAccessControl()),
				apiMiddleware,
				apiMiddleware2,
			)
			api.ServeHTTP(w, r)
		default:
			switch r.Method {
			case "GET":
				web := NewWebsite(l, v, NewAccessControl(), a.config.Website.Assets, a.config.Website.PublicFolder)
				web.ServeHTTP(w, r)
			default:
				ren.ErrMethodNotAllowed()(w, r)
			}
		}
	}
}
