package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"strconv"

	pb "github.com/chukmunnlee/gobgg/messages"

	"google.golang.org/grpc"
)

const SERVER = "localhost"
const PORT = 50051

func checkError(msg string, err error) {
	if nil != err {
		log.Fatalf("%s: %v\n", msg, err)
	}
}

func main() {

	if len(os.Args) <= 1 {
		log.Fatalln("Please provide the game id\n")
	}

	id, err := strconv.ParseInt(os.Args[1], 10, 64)
	checkError(fmt.Sprintf("id error: %s", os.Args[1]), err)
	fmt.Printf("Find game by id: %d", id)

	endpoint := fmt.Sprintf(fmt.Sprintf("%s:%d", SERVER, PORT))

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	checkError(fmt.Sprintf("Cannot connect to %s", endpoint), err)
	defer conn.Close()

	client := pb.NewBoardgamesGeekServiceClient(conn)

	req := &pb.FindGameByIdRequest{
		Id: uint64(id),
	}

	result, err := client.FindGameById(context.Background(), req)
	checkError("Error calling FindGameById", err)

	fmt.Printf("\nresult\n:%v\n", result)
}
