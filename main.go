package main

import (
	"net/http"
	"os"

	sv "github.com/endo-checker/common/server"
	"github.com/joho/godotenv"

	"github.com/endo-checker/logger/handler"
	pbcnn "github.com/endo-checker/logger/internal/gen/logger/v1/loggerv1connect"
	"github.com/endo-checker/logger/store"
)

type Server struct {
	*http.ServeMux
}

func main() {
	godotenv.Load()
	port := ":" + os.Getenv("PORT")
	uri := os.Getenv("MONGO_URI")

	svc := &handler.LoggerServer{
		Store: store.New(uri),
	}
	path, hndlr := pbcnn.NewLoggerServiceHandler(svc)

	srvr := sv.Server{
		ServeMux: &http.ServeMux{},
	}

	sv.Server.ConnectServer(srvr, path, hndlr, port)
}