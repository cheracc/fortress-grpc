package handlers

import (
	"net"

	"github.com/cheracc/fortress-grpc"
	fgrpc "github.com/cheracc/fortress-grpc/grpc"
	"google.golang.org/grpc"
)

// A GrpcServer handles all gRPC requests. It contains the http server and listener used as well as a reference to the Logger
type GrpcServer struct {
	*grpc.Server
	net.Listener
	*fortress.Logger
}

// NewGrpcServer constructs a new GrpcServer with the given logger
func NewGrpcServer(logger *fortress.Logger) GrpcServer {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalf("failed to start a tcp listener on port 50051: %v", err)
	}

	server := grpc.NewServer()

	return GrpcServer{server, listener, logger}
}

func (h *GrpcServer) RegisterHandlers(p *PlayerHandler, a *AuthHandler, cmd *CommandHandler, chat *ChatHandler) {
	fgrpc.RegisterAuthServer(h.Server, a)
	fgrpc.RegisterCommandServer(h.Server, cmd)
	fgrpc.RegisterPlayerServer(h.Server, p)
	fgrpc.RegisterChatServer(h.Server, chat)
}

// StartListener starts the server. It must only be called once all other receivers have been registered
func (h *GrpcServer) StartListener() {
	h.Log("Starting gRPC Server")
	err := h.Server.Serve(h.Listener)
	if err != nil {
		h.Logf("failed to serve: %v", err)
	}
	h.Logf("gRCP listening on %s", h.Listener.Addr().String())
}
