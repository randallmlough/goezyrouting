package website

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
)

func NewViews(templateDir, baseTemplate, viewsDir, siteTitle string, styles []string) (*View, error) {
	t, err := NewTemplates().ParseDir("./"+templateDir, templateDir+"/")
	if err != nil {
		return nil, errors.New("failed to parse template directory")
	}

	t.AddFuncs(
		FuncMap,
	)

	t.Parse()
	v := &View{
		Site: &Site{
			Title:  siteTitle,
			Styles: styles,
			Page:   &Page{},
		},
		baseTemplate: baseTemplate,
		viewsDir:     viewsDir,
		templates:    t,
	}

	return v, nil
}

type View struct {
	Site         *Site
	baseTemplate string
	viewsDir     string
	templates    *Templates
}

type Site struct {
	Title  string
	Styles []string
	Page   *Page
}

type Page struct {
	URI   string
	Error error
	Data  interface{}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	switch d := data.(type) {
	case *Page:
		v.Site.Page = d
	case Page:
		v.Site.Page = &d
	default:
		v.Site.Page = &Page{
			Data: data,
		}
	}

	v.Site.Page.URI = r.URL.Path
	b, err := v.templates.Render(v.baseTemplate, v.viewsDir+view, v.Site)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	// If we get here that means our template executed correctly
	// and we can copy the buffer to the ResponseWriter
	if _, err := io.Copy(w, bytes.NewReader(b)); err != nil {
		return err
	}
	return nil
}
