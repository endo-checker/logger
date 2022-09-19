package main

import (
	"context"
	"log"
	"net/http"

	gw "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/endo-checker/logger/gen/proto/go/logger/v1"
	"github.com/endo-checker/logger/handler"
	"github.com/endo-checker/logger/store"

	sv "github.com/endo-checker/common/server"
)

const defPort = ":8081"

func main() {
	grpcSrv := grpc.NewServer()
	defer grpcSrv.Stop()         // stop server on exit
	reflection.Register(grpcSrv) // for postman

	h := &handler.LoggerServer{
		Store: store.Connect(),
	}

	hm := gw.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		switch key {
		case "X-Token-C-Tenant", "X-Token-C-User", "Permissions":
			return key, true
		default:
			return gw.DefaultHeaderMatcher(key)
		}
	})
	mo := gw.WithMarshalerOption("*", &gw.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: false,
		},
	})

	pb.RegisterLoggerServiceServer(grpcSrv, h)
	httpMux := gw.NewServeMux(hm, mo)

	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterLoggerServiceHandlerFromEndpoint(context.Background(), httpMux, defPort, dopts); err != nil {
		log.Fatal(err)
	}
	mux := sv.HttpGrpcMux(httpMux, grpcSrv)
	httpSrv := &http.Server{
		Addr:    defPort,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		return
	}
}
