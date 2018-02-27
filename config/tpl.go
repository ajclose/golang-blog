package config

import (
	"html"
	"html/template"
)

var TPL *template.Template

func init() {
	TPL = template.Must(template.New("").Funcs(fm).ParseGlob("views/*"))
}

func unescape(s string) string {
	return html.UnescapeString(s)
}

var fm = template.FuncMap{
	"unescape": unescape,
}
