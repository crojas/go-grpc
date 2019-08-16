package main

import (
	"fmt"
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
	fmt.Printf("Cliente creado %f", c)
}
