package main

import (
	"fmt"
	"net/http"
)

func main() {
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	logger.Info("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("Error starting the server: " + err.Error())
	}
}
