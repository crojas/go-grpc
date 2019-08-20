package main

import (
	"context"
	"fmt"
	"log"

	"calculator-grpc/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Cliente GRPC calculator")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("erro al conectar %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 11,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("erro al llamar a servicio SUM GRPC %v", err)
	}
	log.Printf("Respuesta SUM GRPC: %v", res.SumResult)
}
