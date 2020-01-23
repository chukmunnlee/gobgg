package main

import (
	"fmt"
	"os"
	"strconv"
)

const MONGO_URL = "mongodb://localhost:27017"

func main() {
	id, err := strconv.ParseInt(os.Args[1], 10, 64)
	if nil != err {
		fmt.Printf("Invalid id: %v\n", err)
		panic(err)
	}

	s := &BggService{
		MongoURL: MONGO_URL,
	}

	if err := s.Connect(); nil != err {
		fmt.Printf("Connection error: %v\n", err)
		panic(err)
	}

	game, err := s.bggFindById(uint64(id))
	if nil != err {
		fmt.Printf("Query error: %v\n", err)
		panic(err)
	}

	fmt.Printf("result: %v\n", game)
}
