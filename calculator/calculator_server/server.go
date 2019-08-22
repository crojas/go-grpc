package main

import (
	"calculator-grpc/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("funcion Sum invocada con req %v\n", req)
	res := req.GetFirstNumber() + req.GetSecondNumber()
	result := &calculatorpb.SumResponse{
		SumResult: res,
	}
	return result, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("funcion PrimeNumberDecomposition invocada con req %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Println("Aumentando el dividor a %v", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("funcion ComputeAverage invocada con stream del cliente")
	sumTotal := 0.0
	average := 0.0
	count := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("error recibiendo numeros %v", err)
		}
		sumTotal += req.GetNumber()
		count++
		average = sumTotal / count
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("funcion FindMaximum invocada con stream del cliente")
	var max int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error recibiendo mensaje %v\n", err)
		}
		if max <= req.GetNumber() {
			max = req.GetNumber()
		}
		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("Error enviando maximo %v\n", err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error al escuchar el puerto 50051 %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir %v", err)
	}
}
