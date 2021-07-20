package main

import (
	"flag"
	"fmt"
	"github.com/sakiib/grpc-try/gen/pb"
	"github.com/sakiib/grpc-try/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()

	log.Printf("starting the server on port: %d", *port)

	grpcServer := grpc.NewServer()
	bookServer := service.NewBookService(service.NewInMemStore())
	reflection.Register(grpcServer)

	pb.RegisterBookServiceServer(grpcServer, bookServer)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to liten with: %s", err.Error())
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve with: %s", err.Error())
	}
}
