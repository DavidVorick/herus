package main

// server.go manages the server that handles requests.

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// rootHandler will handle any request that comes to the page root - usually
// indicating some type of malformed request or 404.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(filepath.Join(dirTemplates, "404.tpl"))
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

// establishServerRoutes writes all of the routes that are understandable to
// the server.
func (h *herus) establishServerRoutes() {
	http.HandleFunc("/", rootHandler)

	// Feed all of the web assets - must be accessible from /t/ and /m/ as
	// well.
	wafs := http.FileServer(http.Dir("web-assets"))
	wafs = http.StripPrefix("/web-assets/", wafs)
	http.Handle("/web-assets/", wafs)

	http.Handle("/m/", http.StripPrefix("/m/", http.FileServer(http.Dir("media")))) // Serve everything directly from the media dir
	http.HandleFunc(connectPage, h.connectHandler)
	http.HandleFunc(elaborationPrefix, h.elaborationHandler)
	http.HandleFunc(indexPage, indexHandler)
	http.HandleFunc(topicPrefix, h.topicHandler)
	http.HandleFunc(uploadPage, h.uploadHandler)
}
