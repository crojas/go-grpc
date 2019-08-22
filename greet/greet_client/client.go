package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"greet-grpc/greetpb"

	"google.golang.org/grpc"
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
	doBiDirectionalStreaming(c)
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
		log.Fatalf("error llamando a LongGreet", err)
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
