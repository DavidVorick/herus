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
	t, err := template.ParseFiles(filepath.Join(templatesDir, "404.tpl"))
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
	http.HandleFunc(connectPrefix, h.connectHandler)
	http.HandleFunc(elaborationPrefix, h.elaborationHandler)
	http.HandleFunc(indexPrefix, indexHandler)
	http.HandleFunc(mediaPrefix, mediaHandler)
	http.HandleFunc(topicPrefix, h.topicHandler)
	http.HandleFunc(uploadPrefix, h.uploadHandler)
}
