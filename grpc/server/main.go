package main

import (
	"context"
	"log"
	"net"

	hwpb "git.yusank.space/yusank/klyn-examp/grpc/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":7711"
)

type server struct{}

func (s *server) SayHello(c context.Context, r *hwpb.HelloRequest) (resp *hwpb.HelloResponse, err error) {
	log.Printf("receive: %v \n", r.Name)
	return &hwpb.HelloResponse{Message: "Hello " + r.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hwpb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
