package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Server *http.Server
}

const (
	ReadTimeout    = 1 * time.Hour
	WriteTimeout   = 1 * time.Hour
	MaxHeaderBytes = 1 << 28
)

func NewServer(handler http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:           fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
			ReadTimeout:    ReadTimeout,
			WriteTimeout:   WriteTimeout,
			MaxHeaderBytes: MaxHeaderBytes,
			Handler:        handler,
		},
	}
}

func (srv *Server) Run() error {
	log.Printf("Server started with addr: %s\n", srv.Server.Addr)
	return srv.Server.ListenAndServe()
}
