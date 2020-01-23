package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/chukmunnlee/mgbgg/messages"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
{
        "_id" : ObjectId("5d68b4d13cc1a2f130b1224d"),
        "ID" : 174430,
        "Name" : "Gloomhaven",
        "Year" : 2017,
        "Rank" : 1,
        "Average" : 8.9,
        "Bayes average" : 8.611,
        "Users rated" : 24574,
        "URL" : "/boardgame/174430/gloomhaven",
        "Thumbnail" : "https://cf.geekdo-images.com/micro/img/8JYMPXdcBg_UHddwzq64H4NBduY=/fit-in/64x64/pic2437871.jpg"
}
*/

type Game struct {
	_id          primitive.ObjectID `bson:"_id,omitempty"`
	Id           uint64             `bson:"ID"`
	Name         string             `bson:"Name"`
	Year         uint64             `bson:"Year"`
	Rank         uint64             `bson:"Rank"`
	Average      float64            `bson:"Average"`
	BayesAverage float64            `bson:"Bayes average"`
	UsersRated   uint64             `bson:"Users rated"`
	Url          string             `bson:"URL"`
	Thumbnail    string             `bson:"Thumbnail"`
}

type BggService struct {
	MongoURL string
	Client   *mongo.Client
	Games    *mongo.Collection
	pb.UnimplementedBoardgamesGeekServiceServer
}

func toProtobuf(g *Game) *pb.Game {
	if nil != g {
		return &pb.Game{
			Id:           g.Id,
			Name:         g.Name,
			Year:         g.Year,
			Rank:         g.Rank,
			Average:      g.Average,
			BayesAverage: g.BayesAverage,
			UsersRated:   g.UsersRated,
			Url:          g.Url,
			Thumbnail:    g.Thumbnail,
		}
	}
	return nil
}

// gRPC
func (s *BggService) FindGameById(ctx context.Context, req *pb.FindGameByIdRequest) (*pb.FindGameByIdResponse, error) {
	id := req.GetId()
	game, err := s.bggFindById(id)
	if nil != err {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("bggFindById(%f): error = %v", id, err),
		)
	}
	resp := &pb.FindGameByIdResponse{
		Id:   id,
		Game: toProtobuf(game),
	}
	if nil == game {
		resp.Status = pb.Status_NOT_FOUND
	} else {
		resp.Status = pb.Status_OK
	}
	return resp, nil
}

// Mongo utils
func (s *BggService) bggFindById(id uint64) (*Game, error) {
	filter := bson.D{{"ID", id}}
	game := &Game{}
	result := s.Games.FindOne(context.TODO(), filter)
	if nil != result.Err() {
		return nil, nil
	}
	if err := result.Decode(game); nil != err {
		return nil, err
	}
	return game, nil
}

func (s *BggService) Connect() error {
	log.Printf("Connecting to Mongo")

	client, err := mongo.NewClient(options.Client().ApplyURI(s.MongoURL))
	if nil != err {
		return fmt.Errorf("Cannot open connection: %s: %v", s.MongoURL, err)
	}

	if err := client.Connect(context.TODO()); nil != err {
		return fmt.Errorf("Cannot connect to mongo: %s: %v", s.MongoURL, err)
	}
	s.Client = client
	s.Games = client.Database("bgg").Collection("games")

	return nil
}

func (s *BggService) Close() {
	s.Client.Disconnect(context.TODO())
}
