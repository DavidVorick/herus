package main

// server.go manages the server that handles requests.

import (
	"io"
	"net/http"
)

// rootHandler will handle any request that comes to the page root - usually
// indicating some type of malformed request or 404.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "404")
}

// establishServerRoutes writes all of the routes that are understandable to
// the server.
func (h *herus) establishServerRoutes() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc(connectPrefix, h.connectHandler)
	http.HandleFunc(indexPrefix, indexHandler)
	http.HandleFunc(mediaPrefix, mediaHandler)
	http.HandleFunc(topicPrefix, h.topicHandler)
	http.HandleFunc(uploadPrefix, h.uploadHandler)
}
