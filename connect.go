package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	connectPage = "/connect.go"

	connectTitle = "Connect one topic to another"
)

var (
	connectTpl = filepath.Join(dirTemplates, "connect.tpl")

	errDuplicateRelation       = errors.New("relation already exists")
	errMissingSourceTopic      = errors.New("source topic does not exist - cannot add connection")
	errMissingDestinationTopic = errors.New("destination topic does not exist - cannot add connection")
)

// ConnectTemplateData defines the data which is used to fill out the connect
// template file.
type ConnectTemplateData struct {
	Error            string
	ErrorExists      bool
	PostWithoutError bool
}

// connectHandler handles requests to connect pages.
func (h *herus) connectHandler(w http.ResponseWriter, r *http.Request) {
	var ctd ConnectTemplateData
	var err error
	if r.Method == "POST" {
		err = h.processConnectSubmission(r)
		if err != nil {
			ctd.ErrorExists = true
			ctd.Error = err.Error()
		} else {
			ctd.PostWithoutError = true
		}
	}

	// Regardless of the method type, serve the connect page.
	err = executeHeader(w, HeaderTemplateData{Title: connectTitle})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeConnectBody(w, ctd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeFooter(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// executeConnectBody builds the body portion of the connect page.
func executeConnectBody(w io.Writer, ctd ConnectTemplateData) error {
	t, err := template.ParseFiles(connectTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, ctd)
}

// processConnectSubmission processes a request to connect two topics.
func (h *herus) processConnectSubmission(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	sourceTopic := r.FormValue("sourceTopic")
	sourceTopic = strings.Replace(sourceTopic, " ", "_", -1)
	sourceTopic = strings.ToLower(sourceTopic)
	destinationTopic := r.FormValue("destinationTopic")
	destinationTopic = strings.Replace(destinationTopic, " ", "_", -1)
	destinationTopic = strings.ToLower(destinationTopic)

	// Add the topic relation to the source topic db file, but only if the
	// destination topic already exists.
	return h.db.Update(func(tx *bolt.Tx) error {
		sourceTD, exists1, err := getTopicData(tx, sourceTopic)
		if err != nil {
			return err
		}
		_, exists2, err := getTopicData(tx, destinationTopic)
		if err != nil {
			return err
		}
		if !exists1 {
			return errMissingSourceTopic
		}
		if !exists2 {
			return errMissingDestinationTopic
		}

		// Check if the relation being created already exists.
		for _, relation := range sourceTD.RelatedTopics {
			if relation.Title == destinationTopic {
				return errDuplicateRelation
			}
		}

		// Add the relation.
		sourceTD.RelatedTopics = append(sourceTD.RelatedTopics, topicRelation{
			Title:          destinationTopic,
			SubmissionDate: time.Now(),
			// Submitter:

			Downvotes:  0,
			LeftVotes:  0,
			RightVotes: 0,
			Upvotes:    3,
		})
		return putTopicData(tx, sourceTopic, sourceTD)
	})
}
