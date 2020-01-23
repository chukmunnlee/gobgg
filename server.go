package main

import (
	"fmt"
	"log"
	"os"

	"net"

	pb "github.com/chukmunnlee/gobgg/messages"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PROTOCOL = "tcp"
const PORT = 50051
const MONGO_URL = "mongodb://localhost:27017"

func checkError(msg string, err error) {
	if nil != err {
		log.Fatalf("%s: %v\n", msg, err)
	}
}

func main() {

	fmt.Printf("%v\n", os.Args)

	lis, err := net.Listen(PROTOCOL, fmt.Sprintf(":%d", PORT))
	checkError(fmt.Sprintf("Cannot open port: %d", PORT), err)

	s := grpc.NewServer()

	bggSvc := &BggService{
		MongoURL: MONGO_URL,
	}
	checkError("Cannnot open Mongo", bggSvc.Connect())

	reflection.Register(s)

	pb.RegisterBoardgamesGeekServiceServer(s, bggSvc)

	log.Printf("Starting BoardgamesGeekService")
	checkError("Cannot start BoardgamesGeekService", s.Serve(lis))
}
