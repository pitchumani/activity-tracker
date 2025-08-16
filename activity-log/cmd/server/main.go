package main

import (
	"log"
	"net"

	"github.com/pitchumani/activity-tracker/activity-log/internal/server"
	"google.golang.org/grpc/reflection"
)

func main() {
	println("Starting server. Listening on port 8080")
	port := ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	log.Printf("Listening on %s\n", port)
	srv := server.NewGRPCServer()

	// enable reflection service on gRPC servicr
	reflection.Register(srv)

	if err := srv.Serve(lis); err != nil {
		log.Fatal("Failed to serve: %v\n", err)
	}
}
