package main

import (
	"io"
	"net/http"
)

const (
	indexPrefix = "/index"
)

// indexHandler will handle any requests coming to the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Index page of knosys")
}
