package main

// media.go manages the information regarding a piece of media, including
// elaborations and annotations that have been added to the media.

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
	elaborationPrefix = "/e/"
	mediaPrefix       = "/m/"
)

var (
	elaborationTpl = filepath.Join(dirTemplates, "elaborations.tpl")
)

// mediaMetadata contains the metadata regarding a piece of media.
type mediaMetadata struct {
	Title        string
	Elaborations []mediaElaboration
}

// elaborationTemplateData creates the data that fills out
// templates/elaborations.tpl
type elaborationTemplateData struct {
	ElaborationPrefix string
	Hash              string
	MediaPrefix       string
	Title             string

	Elaborations []mediaElaboration
}

// mediaElaboration connects elaborating media to the source media.
type mediaElaboration struct {
	Hash           string
	SubmissionDate time.Time
	Submitter      string
	Title          string

	Downvotes uint64
	Upvotes   uint64
}

// executeMediaBody writes the html body for the media page.
func executeMediaBody(w io.Writer, etd elaborationTemplateData) error {
	t, err := template.ParseFiles(elaborationTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, etd)
}

// getMediaMetadata pulls the elaborations on a piece of media from the
// database.
func getMediaMetadata(tx *bolt.Tx, mediaName string) (mm mediaMetadata, exists bool, err error) {
	// Get the data from the bucket.
	bm := tx.Bucket(bucketMedia)
	mmBytes := bm.Get([]byte(mediaName))
	if mmBytes == nil {
		return mediaMetadata{}, false, nil
	}

	// Unmarshal the data.
	err = json.Unmarshal(mmBytes, &mm)
	if err != nil {
		return mediaMetadata{}, true, err
	}
	return mm, true, nil
}

// elaborationHandler handles requests for the elaborations on media.
func (h *herus) elaborationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the hash of the media for which the elaboration is being loaded.
	sourceHash := strings.TrimPrefix(r.URL.Path, elaborationPrefix)

	// Get a list of elaborations on the source media.
	var mm mediaMetadata
	err := h.db.View(func(tx *bolt.Tx) error {
		var exists bool
		var err error
		mm, exists, err = getMediaMetadata(tx, sourceHash)
		if !exists {
			return nil
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fill out the struct that informs the elaboration template.
	etd := elaborationTemplateData{
		ElaborationPrefix: elaborationPrefix,
		Hash:              sourceHash,
		MediaPrefix:       mediaPrefix,

		Elaborations: mm.Elaborations,
	}
	err = executeHeader(w, HeaderTemplateData{Title: topicTitle})
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeMediaBody(w, etd)
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
