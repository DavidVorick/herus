package main

import (
	"io"
	"net/http"
	"strings"
)

const (
	topicPrefix = "/t/"
)

// topicHandler handles requests for topic pages.
func topicHandler(w http.ResponseWriter, r *http.Request) {
	desiredPage := strings.TrimPrefix(r.URL.Path, topicPrefix)
	io.WriteString(w, "got a topic page: " + desiredPage)
}