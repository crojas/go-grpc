package main

import (
	"context"
	"fmt"
	"io"
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
	//doUnary(c)
	doServerStreaming(c)
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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 210,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error llamando a la funcion PrimeNumberDecomposition del grpc %v", err)
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
		log.Printf("Respuesta PrimeNumberDecomposition %v", msg.GetPrimeFactor())
	}
}
