package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	doServerStreaming(c)
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
