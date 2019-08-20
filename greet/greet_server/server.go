package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"greet-grpc/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

// Unary call
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("funcion Greet invocada con req %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hola " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

// Sever streaming call
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("funcion GreetManyTimes invocada con req %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hola " + firstName + " numero " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("Mi primer GRPC")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error al escuchar en el puerto 50051 %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir %v", err)
	}
}
