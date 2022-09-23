package handler

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/endo-checker/logger/gen/proto/go/logger/v1"
	"github.com/endo-checker/logger/store"
)

type LoggerServer struct {
	Store store.Storer
	Dapr  dapr.Client
	pb.UnimplementedLoggerServiceServer
}

func (log LoggerServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	logs := req.Log
	logs.Id = uuid.NewString()
	logs.Date = time.Now().Unix()

	if err := log.Store.AddLog(logs, md); err != nil {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}
	return &pb.CreateResponse{Log: logs}, nil
}

func (log LoggerServer) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.QueryResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	cur, mat, err := log.Store.QueryLog(req, md)
	if err != nil {
		return &pb.QueryResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	return &pb.QueryResponse{
		Cursor:  cur,
		Matches: mat,
	}, nil
}

func (log LoggerServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.GetResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	ptnt, err := log.Store.GetLog(req.Id, md)
	if err != nil {
		return &pb.GetResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	return &pb.GetResponse{Log: ptnt}, nil
}

func (log LoggerServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.UpdateResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	l := req.Log
	id := l.Id

	if err := log.Store.UpdateLog(id, md, l); err != nil {
		return &pb.UpdateResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	return &pb.UpdateResponse{Log: l}, nil
}

func (log LoggerServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.DeleteResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	if err := log.Store.DeleteLog(req.Id, md); err != nil {
		return &pb.DeleteResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	return &pb.DeleteResponse{}, nil
}
