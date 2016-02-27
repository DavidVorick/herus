package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	mediaPrefix = "/media/"
)

// mediaMetadata contains information about media that has been submitted to
// herus.
type mediaMetadata struct {
	Hash           string
	SubmissionDate time.Time
	Submitter      string
	Title          string

	Downvotes  uint64
	LeftVotes  uint64
	RightVotes uint64
	Upvotes    uint64
}

// mediaHandler will serve the media found with the given url.
func mediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaLocation := strings.TrimPrefix(r.URL.Path, mediaPrefix)
	media, err := ioutil.ReadFile(filepath.Join(mediaDir, mediaLocation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write(media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}
