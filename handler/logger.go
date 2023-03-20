package handler

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/bufbuild/connect-go"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/google/uuid"

	pb "github.com/endo-checker/logger/internal/gen/logger/v1"
	pbcnn "github.com/endo-checker/logger/internal/gen/logger/v1/loggerv1connect"
	"github.com/endo-checker/logger/store"
)

type LoggerServer struct {
	Store store.Storer
	Dapr  dapr.Client
	pbcnn.UnimplementedLoggerServiceHandler
}

func (log LoggerServer) Create(ctx context.Context, req *connect.Request[pb.CreateRequest]) (*connect.Response[pb.CreateResponse], error) {
	reqMsg := req.Msg

	logs := reqMsg.Log
	logs.Id = uuid.NewString()
	logs.Date = time.Now().Unix()

	if err := log.Store.Create(ctx, logs); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.CreateResponse{
		Log: logs,
	}
	return connect.NewResponse(rsp), nil
}

func (log LoggerServer) Query(ctx context.Context, req *connect.Request[pb.QueryRequest]) (*connect.Response[pb.QueryResponse], error) {
	reqMsg := req.Msg

	if reqMsg.LogId != "" {
		pattern, err := regexp.Compile(`^[a-zA-Z0-9-]+$`)
		if err != nil {
			return nil, connect.NewError(connect.CodeAborted, err)
		}
		if !pattern.MatchString(reqMsg.LogId) {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				errors.New("invalid search text format"))
		}
	}

	cur, mat, err := log.Store.Fetch(ctx, reqMsg)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.QueryResponse{
		Cursor:  cur,
		Matches: mat,
	}
	return connect.NewResponse(rsp), nil
}

func (log LoggerServer) Get(ctx context.Context, req *connect.Request[pb.GetRequest]) (*connect.Response[pb.GetResponse], error) {
	reqMsg := req.Msg

	l, err := log.Store.Get(ctx, reqMsg.LogId)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.GetResponse{Log: l}
	return connect.NewResponse(rsp), nil
}

func (log LoggerServer) Update(ctx context.Context, req *connect.Request[pb.UpdateRequest]) (*connect.Response[pb.UpdateResponse], error) {
	reqMsg := req.Msg

	if err := log.Store.Update(ctx, reqMsg.LogId, reqMsg.Log); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.UpdateResponse{
		Log: reqMsg.Log,
	}
	return connect.NewResponse(rsp), nil

}

func (log LoggerServer) Delete(ctx context.Context, req *connect.Request[pb.DeleteRequest]) (*connect.Response[pb.DeleteResponse], error) {
	reqMsg := req.Msg

	if err := log.Store.Delete(ctx, reqMsg.LogId); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.DeleteResponse{}
	return connect.NewResponse(rsp), nil
}
