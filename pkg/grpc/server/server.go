package server

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// GRPCServer GRPCServer
type GRPCServer struct {
	listener net.Listener
	Grpc     *grpc.Server
}

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var keep = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               2 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	log.Printf("request - method: %s\tduration: %s\terror: %v\n",
		info.FullMethod,
		time.Since(start),
		err,
	)

	return h, err
}

// NewServer NewServer
func NewServer(port string) *GRPCServer {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil
	}

	creds, err := credentials.NewServerTLSFromFile("./certs/server.crt", "./certs/server.key")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	return &GRPCServer{
		listener: listener,
		Grpc: grpc.NewServer(
			grpc.Creds(creds),
			grpc.UnaryInterceptor(unaryInterceptor),
			grpc.KeepaliveEnforcementPolicy(kaep),
			grpc.KeepaliveParams(keep),
		),
	}
}

// Start Start
func (r *GRPCServer) Start() error {
	return r.Grpc.Serve(r.listener)
}

// Stop Stop
func (r *GRPCServer) Stop() {
	r.Grpc.GracefulStop()
}
