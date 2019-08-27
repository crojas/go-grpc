package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/crojas/go-grpc-course/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) GreetWithDeadLine(ctx context.Context, req *greetpb.GreetWithDeadLineRequest) (*greetpb.GreetWithDeadLineResponse, error) {
	fmt.Printf("funcion GreetWithDeadLine invocada con cliente v%", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// El cliente cancelo la peticion
			fmt.Println("El cliente ha cancelado la peticion")
			return nil, status.Error(codes.Canceled, "el cliente ha cancelado la peticion")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hola " + firstName
	res := &greetpb.GreetWithDeadLineResponse{
		Result: result,
	}
	return res, nil
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
