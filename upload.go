package main

// upload.go handles the uploading of data to the server.

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const (
	uploadPrefix = "/upload"
)

// receiveUpload accepts an upload presented by the user.
func receiveUpload(w http.ResponseWriter, r *http.Request) {
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
	filename := hex.EncodeToString(checksum) + ".txt" // TODO: Accept '.txt', '.png', and '.pdf'.
	err = ioutil.WriteFile(filepath.Join(mediaDir, filename), fileData, 0700)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := template.ParseFiles("uploadSuccess.gtpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, nil)
}

// serveUploadPage presents the page that users can use to upload files to the
// server.
func serveUploadPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("upload.gtpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, nil)
}

// uploadHandler handles requests for the upload page.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" || r.Method == "GET" {
		serveUploadPage(w, r)
		return
	}
	if r.Method == "POST" {
		receiveUpload(w, r)
		return
	}
}
