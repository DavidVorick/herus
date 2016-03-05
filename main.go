package main

// main.go intializes the herus server and all of its components.

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	dbFile       = "herus.db"
	mediaDir     = "media"
	templatesDir = "templates"
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
			bucektUsers,
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
	err = os.MkdirAll(mediaDir, 0700)
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
