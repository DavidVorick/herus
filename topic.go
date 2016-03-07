package main

// topic.go produces the topics page.

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	topicPrefix = "/t/"

	topicTitle = "Knowledge"
)

var (
	topicTpl = filepath.Join(dirTemplates, "topic.tpl")
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
	Title             string

	AssociatedMedia []mediaRelation
	RelatedTopics   []topicRelation
}

// executeTopicBody writes the html body for the topic page.
func executeTopicBody(w io.Writer, ttd topicTemplateData) error {
	t, err := template.ParseFiles(topicTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, ttd)
}

// getTopic returns the topic data associated with a topic in the bucket.
func getTopic(tx *bolt.Tx, topic string) (td topicData, exists bool, err error) {
	// Get the topic data, checking whether the data exists.
	bt := tx.Bucket(bucketTopics)
	topicDataBytes := bt.Get([]byte(topic))
	if topicDataBytes == nil {
		return topicData{}, false, nil
	}

	// Unmarshal the topic data
	err = json.Unmarshal(topicDataBytes, &td)
	if err != nil {
		return topicData{}, true, err
	}
	return td, true, nil
}

// putTopic stores the provided topic data in the topic database.
func putTopic(tx *bolt.Tx, topic string, td topicData) error {
	topicDataBytes, err := json.Marshal(td)
	if err != nil {
		return err
	}
	bt := tx.Bucket(bucketTopics)
	return bt.Put([]byte(topic), topicDataBytes)
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
	var exists bool
	var err error
	err = h.db.View(func(tx *bolt.Tx) error {
		td, exists, err = getTopic(tx, topicName)
		if !exists {
			return nil
		}
		return err
	})
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fill out the stuct that will inform the topic template.
	ttd := topicTemplateData{
		ElaborationPrefix: elaborationPrefix,
		Title:             topicTitle,

		AssociatedMedia: td.AssociatedMedia,
		RelatedTopics:   td.RelatedTopics,
	}

	// Execute a template to display all of the uploaded media.
	err = executeHeader(w, HeaderTemplateData{Title: topicTitle})
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeTopicBody(w, ttd)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeFooter(w)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
