package main

import (
	"context"
	"fmt"
	"io"
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

// Client streaming call
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("funcion LongGreet invocada con stream del cliente")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error al recibir mensajes del cliente")
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hola " + firstName + "! "
	}
}

// Bi directional streaming call
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("funcion GreetEveryone invocada con stream del cliente")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error recibiendo stream del cliente %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hola " + firstName + " !"
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error enviando data al cliente %v", err)
		}
	}
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
