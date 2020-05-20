package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// GRPCServer GRPCServer
type GRPCServer struct {
	listener net.Listener
	Grpc     *grpc.Server
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

var errInvalidToken = errors.New("Invalid token")

func jwtUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, _ := metadata.FromIncomingContext(ctx)

	tokenHeader := meta.Get("authorization")
	if len(tokenHeader) == 0 {
		return nil, errInvalidToken
	}

	jwtToken := tokenHeader[0]
	if len(jwtToken) == 0 {
		return nil, errInvalidToken
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("APP_JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, err
	}

	h, err := handler(ctx, req)

	return h, err
}

func grpcMiddlewares() []grpc.UnaryServerInterceptor {
	unary := []grpc.UnaryServerInterceptor{
		unaryInterceptor,
	}

	if useJWT := os.Getenv("APP_JWT"); useJWT == "1" {
		log.Println("Middleware: JWT")

		jwtSecret := os.Getenv("APP_JWT_SECRET")
		if len(jwtSecret) == 0 {
			log.Fatalln("Invalid jwt secret", jwtSecret)
		}

		unary = append(unary, jwtUnaryInterceptor)
	}

	return unary
}

// NewServer NewServer
func NewServer(port string) *GRPCServer {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	creds, err := credentials.NewServerTLSFromFile("./certs/server.crt", "./certs/server.key")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	return &GRPCServer{
		listener: listener,
		Grpc: grpc.NewServer(
			grpc.Creds(creds),
			grpc.ChainUnaryInterceptor(grpcMiddlewares()...),
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
	r.listener.Close()
}
