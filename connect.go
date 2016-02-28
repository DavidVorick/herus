package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	connectPrefix = "/connect"
)

var (
	errDuplicateRelation = errors.New("relation already exists!")
	errMissingTopic      = errors.New("either the source or destination topic does not exist - cannot add connection")
)

// receiveConnect handles a connection post request.
func (h *herus) receiveConnect(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sourceTopic := r.FormValue("sourceTopic")
	sourceTopic = strings.Replace(sourceTopic, " ", "_", -1)
	sourceTopic = strings.ToLower(sourceTopic)
	destinationTopic := r.FormValue("destinationTopic")
	destinationTopic = strings.Replace(destinationTopic, " ", "_", -1)
	destinationTopic = strings.ToLower(destinationTopic)

	// Add the topic relation to the source topic db file, but only if the
	// destination topic already exists.
	err = h.db.Update(func(tx *bolt.Tx) error {
		// First check that both the source and destination topics exist.
		bt := tx.Bucket(bucketTopics)
		sourceTopicDataBytes := bt.Get([]byte(sourceTopic))
		destData := bt.Get([]byte(destinationTopic))
		if sourceTopicDataBytes == nil || destData == nil {
			return errMissingTopic
		}

		// Get the topic data.
		var td topicData
		err = json.Unmarshal(sourceTopicDataBytes, &td)
		if err != nil {
			return nil
		}

		// Check if the relation being created already exists.
		for _, relation := range td.RelatedTopics {
			if relation.Title == destinationTopic {
				return errDuplicateRelation
			}
		}

		// Add the relation.
		td.RelatedTopics = append(td.RelatedTopics, topicRelation{
			Title:          destinationTopic,
			SubmissionDate: time.Now(),
			// Submitter:

			Downvotes:  0,
			LeftVotes:  0,
			RightVotes: 0,
			Upvotes:    3,
		})
		sourceTopicDataBytes, err = json.Marshal(td)
		if err != nil {
			return err
		}
		return bt.Put([]byte(sourceTopic), sourceTopicDataBytes)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := template.ParseFiles("templates/connect.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, true)
}

// serveConnectTopic presents the page that users can use to upload files to the
// server.
func (h *herus) serveConnectTopic(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/connect.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, false)
}

// connectHandler handles requests to connect pages.
func (h *herus) connectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" || r.Method == "GET" {
		h.serveConnectTopic(w, r)
		return
	}
	if r.Method == "POST" {
		h.receiveConnect(w, r)
		return
	}
}
