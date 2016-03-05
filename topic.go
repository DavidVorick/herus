package main

// topic.go produces the topics page.

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	topicPrefix = "/t/"
)

// mediaRelation contains information about media that has been submitted to
// herus.
type mediaRelation struct {
	Hash           string
	SubmissionDate time.Time
	Submitter      string
	Title          string

	Downvotes  uint64
	LeftVotes  uint64
	RightVotes uint64
	Upvotes    uint64
}

// topicRelation is a mapping from one topic to another.
type topicRelation struct {
	SubmissionDate time.Time
	Submitter      string
	Title          string

	CenterVotes uint64
	Downvotes   uint64
	LeftVotes   uint64
	RightVotes  uint64
	Upvotes     uint64
}

// topicData is a struct that gets stored in the database containing all of the
// information about a topic.
type topicData struct {
	AssociatedMedia []mediaRelation
	RelatedTopics   []topicRelation
}

// topicTemplateData provides the dynamic data that is used to fill out the
// template for the topic page.
type topicTemplateData struct {
	ElaborationPrefix string
	TopicPrefix       string
	TopicTitle        string

	AssociatedMedia []mediaRelation
	RelatedTopics   []topicRelation
}

// topicHandler handles requests for topic pages.
func (h *herus) topicHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and preprocess the topic name.
	//
	// TODO: There is probably a better way to format the topic title.
	topicTitle := strings.TrimPrefix(r.URL.Path, topicPrefix)
	topicTitle = strings.Replace(topicTitle, "_", " ", -1)
	topicTitle = strings.Title(topicTitle)
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
		ElaborationPrefix: elaborationPrefix,
		TopicTitle:        topicTitle,

		AssociatedMedia: td.AssociatedMedia,
		RelatedTopics:   td.RelatedTopics,
	}

	// Execute a template to display all of the uploaded media.
	t, err := template.ParseFiles(filepath.Join(dirTemplates, "topic.tpl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.Execute(w, ttd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
