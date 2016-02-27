package main

// topic.go serves pages associated with the topics requested by users.

import (
	"encoding/json"
	"html/template"
	"net/http"
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

type SubmittedMedia struct {
	MediaTitle string
	MediaHash  string
}

type topicTemplateResults struct {
	TopicTitle     string
	SubmittedMedia []SubmittedMedia
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

	// Fill out the stuct that will inform the topic template.
	ttr := topicTemplateResults{
		TopicTitle: topicName,
	}
	for i := range td.MediaTitles {
		ttr.SubmittedMedia = append(ttr.SubmittedMedia, SubmittedMedia{
			MediaTitle: td.MediaTitles[i],
			MediaHash:  mediaPrefix + td.MediaHashes[i],
		})
	}

	// Execute a template to display all of the uploaded media.
	t, err := template.ParseFiles("topic.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, ttr)
}
