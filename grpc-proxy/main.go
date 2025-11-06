package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/pitchumani/activity-tracker/activity-log/api/v1"
)

func main() {
	var grpcServerEndpoint = ":8090"

	log.Println("Listening to port 8081")
	port := ":8081"
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := api.RegisterActivity_LogHandlerFromEndpoint(context.Background(),
		mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
