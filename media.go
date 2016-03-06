package main

// media.go manages the information regarding a piece of media, including
// elaborations and annotations that have been added to the media.

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
	elaborationPrefix = "/e/"
	mediaPrefix       = "/m/"
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

// elaborationHandler handles requests for the elaborations on media.
func (h *herus) elaborationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the hash of the media for which the elaboration is being loaded.
	sourceHash := strings.TrimPrefix(r.URL.Path, elaborationPrefix)

	// Get a list of elaborations on the source media.
	var mm mediaMetadata
	err := h.db.View(func(tx *bolt.Tx) error {
		// Get the elaboration data.
		bm := tx.Bucket(bucketMedia)
		mediaMetadataBytes := bm.Get([]byte(sourceHash))
		if mediaMetadataBytes != nil {
			err := json.Unmarshal(mediaMetadataBytes, &mm)
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
	t, err := template.ParseFiles(filepath.Join(dirTemplates, "elaborations.tpl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.Execute(w, etd)
	if err != nil {
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
