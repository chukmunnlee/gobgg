package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "github.com/chukmunnlee/gobgg/messages"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const PROTOCOL = "tcp"
const GRPC_PORT = 50051
const HTTP_PORT = 8080

func checkError(msg string, err error) {
	if nil != err {
		log.Fatalf("%s: %v\n", msg, err)
	}
}

func main() {

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	svcEndpoint := fmt.Sprintf(":%d", GRPC_PORT)
	err := pb.RegisterBoardgamesGeekServiceHandlerFromEndpoint(context.Background(), mux, svcEndpoint, opts)
	checkError(fmt.Sprintf("Error registering Endpoint on %d -> %s", HTTP_PORT, svcEndpoint), err)

	log.Printf("Starting BoardgamesGeekService Web Proxy on %d\n", HTTP_PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", HTTP_PORT), mux)
	checkError("Cannot start Web Proxy", err)
}
