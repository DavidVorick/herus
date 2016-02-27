package main

// index.go formats and serves the index page for the herus website.

import (
	"io"
	"net/http"
)

const (
	indexPrefix = "/index"
)

// indexHandler will handle any requests coming to the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Index page of herus")
}
