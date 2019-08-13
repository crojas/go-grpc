package main

import (
	"fmt"
	"log"
	"net"

	"greet-grpc/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Mi primer GRPC")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error al escuchar en el puerto 50051", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir", err)
	}
}
