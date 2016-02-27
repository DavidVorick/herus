package main

// upload.go handles the uploading of data to the server.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/boltdb/bolt"
)

const (
	uploadPrefix = "/upload"
)

// receiveUpload accepts an upload presented by the user.
func (h *herus) receiveUpload(w http.ResponseWriter, r *http.Request) {
	// TODO: The number here indicates the maximum amount of memory that the
	// server will use to parse the file. If the memory goes over, a temp file
	// will be used. But, there should be some way to set a limit on the max
	// size allowed, and I don't think that's happening here.
	r.ParseMultipartForm(8 << 20) // 8 MB
	file, _, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Pull the filedata from the handler.
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the hash of the file.
	//
	// TODO: need to check the size of the file here, though ideally the size
	// of the file would be managed mid-upload.
	hasher := sha256.New()
	hasher.Write(fileData)
	checksum := hasher.Sum(nil)

	// Using the hash, save the file to disk.
	mediaHash := hex.EncodeToString(checksum) + ".txt" // TODO: Accept '.txt', '.png', and '.pdf'.
	err = ioutil.WriteFile(filepath.Join(mediaDir, mediaHash), fileData, 0700)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the topic that this upload is being connected to and update the
	// topic's database file to point to the upload.
	topicName := r.FormValue("topic")
	mediaTitle := r.FormValue("title")
	err = h.db.Update(func(tx *bolt.Tx) error {
		// Fetch the existing topic data.
		var td topicData
		tb := tx.Bucket(bucketTopics)
		topicDataBytes := tb.Get([]byte(topicName))
		if topicDataBytes != nil {
			err = json.Unmarshal(topicDataBytes, &td)
			if err != nil {
				return err
			}
		}

		// Add the new link to the topic data.
		td.MediaTitles = append(td.MediaTitles, mediaTitle)
		td.MediaHashes = append(td.MediaHashes, mediaHash)

		// Save the updated topic data.
		topicDataBytes, err = json.Marshal(td)
		if err != nil {
			return err
		}
		return tb.Put([]byte(topicName), topicDataBytes)
	})

	t, err := template.ParseFiles("uploadSuccess.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, nil)
}

// serveUploadPage presents the page that users can use to upload files to the
// server.
func (h *herus) serveUploadPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("upload.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, nil)
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
