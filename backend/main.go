package main

import (
	. "backend/logger"
	"fmt"
	"net/http"
)

func main() {
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	Logger.Info("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
