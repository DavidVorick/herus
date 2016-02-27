package main

import (
	"fmt"
	"net/http"
)

// main initializes the server and then starts serving pages.
func main() {
	fmt.Println("Starting Herus...")
	establishServerRoutes()
	http.ListenAndServe(":3841", nil)
}
