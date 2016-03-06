package main

// main.go intializes the herus server and all of its components.

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

const (
	dbFile = "herus.db"

	dirCSS       = "css"
	dirMedia     = "media"
	dirTemplates = "templates"
)

var (
	// bucketTopics houses information about all of the pages tracked by herus.
	bucketTopics = []byte("BucketTopics")
	bucketMedia  = []byte("BucketMedia")
	bucketUsers  = []byte("BucketUsers")
)

// herus contains all data that needs to persist in memory throughout the life
// of the server.
type herus struct {
	db *bolt.DB
}

// rootHandler will handle any request that comes to the page root - usually
// indicating some type of malformed request or 404.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(filepath.Join(dirTemplates, "404.tpl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// establishServerRoutes writes all of the routes that are understandable to
// the server.
func (h *herus) establishServerRoutes() {
	http.HandleFunc("/", rootHandler)

	// Feed all of the web assets - must be accessible from /t/ and /m/ as
	// well.
	wafs := http.FileServer(http.Dir("web-assets"))
	wafs = http.StripPrefix("/web-assets/", wafs)
	http.Handle("/web-assets/", wafs)

	http.Handle("/m/", http.StripPrefix("/m/", http.FileServer(http.Dir("media")))) // Serve everything directly from the media dir
	http.HandleFunc(connectPage, h.connectHandler)
	http.HandleFunc(elaborationPrefix, h.elaborationHandler)
	http.HandleFunc(indexPage, indexHandler)
	http.HandleFunc(topicPrefix, h.topicHandler)
	http.HandleFunc(uploadPage, h.uploadHandler)
}

// initDB will initialize the database used by herus.
func (h *herus) initDB() (err error) {
	h.db, err = bolt.Open(dbFile, 0700, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}

	return h.db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			bucketTopics,
			bucketMedia,
			bucketUsers,
		}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// main initializes the server and then starts serving pages.
func main() {
	fmt.Println("Preparing Herus...")
	h := new(herus)

	// Initialize the database.
	err := h.initDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the media folder.
	err = os.MkdirAll(dirMedia, 0700)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set up the server routes.
	h.establishServerRoutes()

	fmt.Println("Serving...")
	err = http.ListenAndServe(":3841", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
