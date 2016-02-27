package main

// main.go intializes the herus server and all of its components.

import (
	"fmt"
	"net/http"
	"os"
)

const (
	mediaDir = "media"
)

// main initializes the server and then starts serving pages.
func main() {
	fmt.Println("Preparing Herus...")
	establishServerRoutes()
	err := os.MkdirAll(mediaDir, 0700)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Serving...")
	http.ListenAndServe(":3841", nil)
}
