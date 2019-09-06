package goezyrouting

import (
	"net/http"
	"strings"
)

func NewWebsite(l Logger, r Renderer, ac AccessController, assets, PublicFolder FileSystem) *Website {
	web := &Website{
		Handler: &Handler{
			l:  l,
			r:  r,
			ac: ac,
		},
		Assets:       assets,
		PublicFolder: PublicFolder,
	}
	return web
}

type Website struct {
	*Handler
	Assets       FileSystem
	PublicFolder FileSystem
}

func (h *Website) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		h.r.ErrMethodNotAllowed()(w, r)
	}

	if strings.HasPrefix(r.URL.Path, h.Assets.Path) {
		assetHandler := http.FileServer(http.Dir(h.Assets.Dir))
		assetHandler = http.StripPrefix(h.Assets.Path, assetHandler)
		assetHandler.ServeHTTP(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, h.PublicFolder.Path) {
		assetHandler := http.FileServer(http.Dir(h.PublicFolder.Dir))
		assetHandler = http.StripPrefix("/"+h.PublicFolder.Path+"/", assetHandler)
		assetHandler.ServeHTTP(w, r)
		return
	}

	//if strings.HasPrefix(r.URL.Path, "/favicon.ico") {
	//	//http.FileServer(http.Dir(h.config.Website.Favicon)).ServeHTTP(w,r)
	//	//http.ServeFile(w, r, h.config.Website.Favicon)
	//	return
	//}

	web := NewStaticHandler(h.l, h.r)
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "user":
		uh := Use(
			NewUserHandler(h.Handler),
			minAccessLevel(3),
			levelTwoMiddleware,
		)
		uh.ServeHTTP(w, r)
	case "account":
	case "login":
		web.login()(w, r)
	case "":
		requireUser(web.Home()).ServeHTTP(w, r)
	default:
		h.r.ErrNotFound()(w, r)
	}
}
