package main

// upload.go handles the uploading of data to the server.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	uploadPage = "/upload.go"

	uploadTitle = "Upload Media to Herus"
)

var (
	uploadTpl = filepath.Join(dirTemplates, "upload.tpl")
)

// receiveUpload accepts an upload presented by the user.
func (h *herus) receiveUpload(w http.ResponseWriter, r *http.Request) {
	// Setting MaxBytesReader limits the file upload to 8 MB, the connection
	// will be closed if the limit is exceeded.
	r.Body = http.MaxBytesReader(w, r.Body, 8<<20) // 8 MB
	err := r.ParseMultipartForm(8 << 20)           // 8 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the values of the elaboration and the topic that the media might be
	// getting added to.
	mediaTitle := r.FormValue("title")
	parentMedia := r.FormValue("parentMedia")
	parentTopic := r.FormValue("parentTopic")
	parentTopic = strings.Replace(parentTopic, " ", "_", -1)
	parentTopic = strings.ToLower(parentTopic)
	if mediaTitle == "" {
		http.Error(w, "media must be uploaded with a title", http.StatusBadRequest)
		return
	}
	if parentMedia == "" && parentTopic == "" {
		http.Error(w, "media must be uploaded with a parent", http.StatusBadRequest)
		return
	}
	if !(parentMedia == "" || parentTopic == "") {
		http.Error(w, "only one parent per upload is currently allowed", http.StatusBadRequest)
		return
	}

	// Pull the file data from the form.
	file, _, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the hash of the file.
	hasher := sha256.New()
	_, err = hasher.Write(fileData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	checksum := hasher.Sum(nil)

	// Create/Update the database entry for this media.
	mediaHash := hex.EncodeToString(checksum)
	var mm mediaMetadata
	mediaExists := false
	oldTitle := mediaTitle
	err = h.db.Update(func(tx *bolt.Tx) error {
		// See if the media has already been added to the server.
		bm := tx.Bucket(bucketMedia)
		mediaMetadataBytes := bm.Get([]byte(mediaHash))
		if mediaMetadataBytes == nil {
			// Create an entry for the media.
			mm.Title = mediaTitle
			mediaMetadataBytes, err := json.Marshal(mm)
			if err != nil {
				return err
			}
			err = bm.Put([]byte(mediaHash), mediaMetadataBytes)
			if err != nil {
				return err
			}
		} else {
			// Update the in-memory title with the in-database title.
			err := json.Unmarshal(mediaMetadataBytes, &mm)
			if err != nil {
				return err
			}
			mediaTitle = mm.Title
			mediaExists = true
		}

		// Add the media to either the parent topic or the parent media.
		if parentMedia != "" {
			// Check whether the media already exists as an elaboration to the
			// parent media.
			var parentMM mediaMetadata
			parentMetadataBytes := bm.Get([]byte(parentMedia))
			if parentMetadataBytes == nil {
				return errors.New("parent media does not exist")
			}
			err = json.Unmarshal(parentMetadataBytes, &parentMM)
			if err != nil {
				return err
			}
			for _, elaboration := range parentMM.Elaborations {
				if elaboration.Hash == mediaHash {
					return errors.New("media has already been added to the parent media")
				}
			}

			// Media does not exist in the parent, add it to the parent.
			parentMM.Elaborations = append(parentMM.Elaborations, mediaElaboration{
				Hash:           mediaHash,
				SubmissionDate: time.Now(),
				// Submitter:
				Title: oldTitle,

				Downvotes: 0,
				Upvotes:   3,
			})
			parentMetadataBytes, err = json.Marshal(parentMM)
			if err != nil {
				return err
			}
			return bm.Put([]byte(parentMedia), parentMetadataBytes)
		}
		// Check whether the parent topic already has the media.
		var td topicData
		bt := tx.Bucket(bucketTopics)
		tdBytes := bt.Get([]byte(parentTopic))
		if tdBytes != nil {
			err = json.Unmarshal(tdBytes, &td)
			if err != nil {
				return err
			}
		}
		for _, am := range td.AssociatedMedia {
			if am.Hash == mediaHash {
				return errors.New("media has already been added to the parent topic")
			}
		}

		// Media does not exist in the parent, add it to the parent.
		td.AssociatedMedia = append(td.AssociatedMedia, mediaRelation{
			Hash:           mediaHash,
			SubmissionDate: time.Now(),
			// Submitter:
			Title: oldTitle,

			Downvotes:  0,
			LeftVotes:  0,
			RightVotes: 0,
			Upvotes:    3,
		})
		tdBytes, err = json.Marshal(td)
		if err != nil {
			return err
		}
		return bt.Put([]byte(parentTopic), tdBytes)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write the media to the media folder.
	if !mediaExists {
		err = ioutil.WriteFile(filepath.Join(dirMedia, mediaHash), fileData, 0700)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = executeHeader(w, HeaderTemplateData{Title: uploadTitle})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := template.ParseFiles(filepath.Join(dirTemplates, "upload.tpl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.Execute(w, true)
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

// executeUploadBody builds the body portion of the upload page.
func executeUploadBody(w io.Writer) error {
	t, err := template.ParseFiles(uploadTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}

// serveUploadPage presents the page that users can use to upload files to the
// server.
func (h *herus) serveUploadPage(w http.ResponseWriter, r *http.Request) {
	err := executeHeader(w, HeaderTemplateData{Title: uploadTitle})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = executeUploadBody(w)
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

// uploadHandler handles requests for the upload page.
func (h *herus) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" || r.Method == "GET" {
		h.serveUploadPage(w, r)
		return
	}
	if r.Method == "POST" {
		h.receiveUpload(w, r)
		return
	}
}
