package config

import (
	"html/template"
	"path/filepath"
)

var Base *template.Template

func layoutFiles() []string {
	files, err := filepath.Glob("views/layouts/*.gohtml")
	if err != nil {
		panic(err)
	}
	return files
}

func CreateView(view string) {
	var err error
	files := append(layoutFiles(), "views/"+view)
	Base, err = template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
}
