package main

// index.go formats and serves the index page for the herus website. It also
// contains all of the code that manages the headers and footers for the
// website.

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	indexPrefix = "/index.go"
)

// indexHandler will handle any requests coming to the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(filepath.Join(templatesDir, "index.tpl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
