package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Addr string `env:"GRPC_ADDR"`
	Port int    `env:"GRPC_PORT"`
}

type Server struct {
	addr string

	ln  net.Listener
	srv *grpc.Server
	cfg *Config
}

func New(config *Config) *Server {
	addr := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	srv := grpc.NewServer()

	return &Server{
		addr: addr,
		srv:  srv,
		cfg:  config,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.ln = ln
	return s.srv.Serve(ln)
}

func (s *Server) Stop() error {
	s.srv.Stop()
	return s.ln.Close()
}

func (s *Server) GracefulStop() error {
	s.srv.GracefulStop()
	return s.ln.Close()
}

func (s *Server) Server() *grpc.Server {
	return s.srv
}
