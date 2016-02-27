package main

// topic.go produces the topics page.

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

// topicRelation is a mapping from one topic to another.
type topicRelation struct {
	TopicTitle string

	CenterVotes uint64
	Downvotes   uint64
	LeftVotes   uint64
	RightVotes  uint64
	Upvotes     uint64
}

// topicData is a struct that gets stored in the database containing all of the
// information about a topic.
type topicData struct {
	AssociatedMedia []mediaMetadata
	RelatedTopics   []topicRelation
}

// topicTemplateData provides the dynamic data that is used to fill out the
// template for the topic page.
type topicTemplateData struct {
	MediaPrefix string
	TopicTitle  string

	AssociatedMedia []mediaMetadata
	RelatedTopics   []topicRelation
}

// topicHandler handles requests for topic pages.
func (h *herus) topicHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and preprocess the topic name.
	//
	// TODO: It's not really great... there's probably a better way to format
	// topicTitle.
	topicTitle := strings.TrimPrefix(r.URL.Path, topicPrefix)
	topicTitle = strings.Replace(topicTitle, "_", " ", -1)
	topicTitle = strings.ToTitle(topicTitle)
	topicName := strings.ToLower(topicTitle)
	topicName = strings.Replace(topicName, " ", "_", -1)

	// Get a list of media from the database and build links to each media
	// file.
	var td topicData
	err := h.db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket(bucketTopics)
		topicDataBytes := bt.Get([]byte(topicName))
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
	ttd := topicTemplateData{
		MediaPrefix: mediaPrefix,
		TopicTitle:  topicTitle,

		AssociatedMedia: td.AssociatedMedia,
		RelatedTopics:   td.RelatedTopics,
	}

	// Execute a template to display all of the uploaded media.
	t, err := template.ParseFiles("templates/topic.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, ttd)
}
