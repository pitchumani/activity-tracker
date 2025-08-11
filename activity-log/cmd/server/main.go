package main

import (
	"github.com/pitchumani/activity-tracker/activity-log/internal/server"
	"log"
)

func main() {
	println("Starting server. Listening on port 8080")
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}

