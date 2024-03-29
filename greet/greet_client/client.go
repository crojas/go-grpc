package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/crojas/go-grpc-course/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Cliente GRPC")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No me pude conectar %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Cliente creado %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDirectionalStreaming(c)
	doUnaryWithDeadline(c, 5*time.Second)
	doUnaryWithDeadline(c, 1*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Carlos",
			LastName:  "Rojas",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error llamando a la funcion Greet del grpc %v", err)
	}
	log.Printf("Respuesta Greet GRPC: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Carlos",
			LastName:  "Rojas",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error llamando a la funcion GreetManyTimes del grpc %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// cuando termina el stream
			break
		}
		if err != nil {
			log.Fatalf("Error leyendo el stream %v", err)
		}
		log.Printf("Respuesta GreetManyTimes %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error llamando a LongGreet %v", err)
	}
	for i := 0; i < 10; i++ {
		msg := "Carlos " + strconv.Itoa(i)
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: msg,
				LastName:  "Rojas",
			},
		}
		fmt.Printf("Enviando %v\n", msg)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error recibiendo mensage del server %v", err)
	}
	fmt.Printf("Respuesta LongGreet %v", res.GetResult())
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error llamando a GreetEveryone")
	}
	waitc := make(chan struct{})
	// Goroutine para enviar mensajes
	go func() {
		for i := 0; i < 10; i++ {
			msg := "Carlos " + strconv.Itoa(i)
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: msg,
				},
			}
			fmt.Printf("Enviando %v\n", msg)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// Goroutine para recibir mensajes
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reciviendo mensajes %v", err)
			}
			fmt.Printf("Recibiendo: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Carlos",
			LastName:  "Rojas",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout, se llego al tiempo limite")
			} else {
				fmt.Printf("error inesperado %v", statusErr)
			}
		} else {
			log.Fatalf("Error llamando a GreetWithDeadLine %v\n", err)
		}
		return
	}
	log.Printf("Respuesta de GreetWithDeadLine %v\n", res.GetResult())
}
