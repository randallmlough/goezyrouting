package website

import (
	"html/template"
)

var FuncMap = template.FuncMap{
	"uri": uri,
}

func uri(data *Page) string {
	return data.URI
}
