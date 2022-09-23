package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/metadata"

	pb "github.com/endo-checker/logger/gen/proto/go/logger/v1"
)

type Storer interface {
	AddLog(u *pb.Log, md metadata.MD) error
	QueryLog(qr *pb.QueryRequest, md metadata.MD) ([]*pb.Log, int64, error)
	GetLog(id string, md metadata.MD) (*pb.Log, error)
	UpdateLog(id string, md metadata.MD, l *pb.Log) error
	DeleteLog(id string, md metadata.MD) error
}

func (s Store) AddLog(logs *pb.Log, md metadata.MD) error {
	_, err := s.locaColl.InsertOne(context.Background(), logs)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s Store) QueryLog(qr *pb.QueryRequest, md metadata.MD) ([]*pb.Log, int64, error) {
	filter := bson.M{}

	if qr.PatientId != "" {
		filter = bson.M{"$text": bson.M{"$search": `"` + qr.PatientId + `"`}}
	}

	opt := options.FindOptions{
		Skip:  &qr.Offset,
		Limit: &qr.Limit,
		Sort:  bson.M{"date": -1},
	}

	ctx := context.Background()
	cursor, err := s.locaColl.Find(ctx, filter, &opt)
	if err != nil {
		return nil, 0, err
	}

	var logs []*pb.Log
	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, 0, err
	}

	matches, err := s.locaColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return logs, matches, err
}

func (s Store) GetLog(id string, md metadata.MD) (*pb.Log, error) {
	var l pb.Log

	if err := s.locaColl.FindOne(context.Background(), bson.M{"id": id}).Decode(&l); err != nil {
		if err == mongo.ErrNoDocuments {
			return &l, err
		}
		return &l, err
	}

	return &l, nil
}

func (s Store) UpdateLog(id string, md metadata.MD, l *pb.Log) error {
	insertResult, err := s.locaColl.ReplaceOne(context.Background(), bson.M{"id": id}, l)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted a Single Document: %v\n", insertResult)

	return err
}

func (s Store) DeleteLog(id string, md metadata.MD) error {
	if _, err := s.locaColl.DeleteOne(context.Background(), bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}
