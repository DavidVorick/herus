package main

// topic.go serves pages associated with the topics requested by users.

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
)

const (
	topicPrefix = "/t/"
)

// topicData is a struct that gets stored in the database containing all of the
// information about a topic.
type topicData struct {
	MediaTitles []string
	MediaHashes []string
}

// topicHandler handles requests for topic pages.
func (h *herus) topicHandler(w http.ResponseWriter, r *http.Request) {
	topicName := strings.TrimPrefix(r.URL.Path, topicPrefix)

	// Get a list of media from the database and build links to each media
	// file.
	var td topicData
	err := h.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket(bucketTopics)
		topicDataBytes := tb.Get([]byte(topicName))
		if topicDataBytes != nil {
			err := json.Unmarshal(topicDataBytes, &td)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write a link for each page.
	io.WriteString(w, "<html><head></head><body><center><h1>"+topicName+"</h1></center>"+"\n")
	for i := range td.MediaTitles {
		mediaLocation := filepath.Join(mediaDir, td.MediaHashes[i])
		io.WriteString(w, "<a href='../"+mediaLocation+"'>"+td.MediaTitles[i]+"</a><br>\n")
	}
	io.WriteString(w, "</body><html>")
}
