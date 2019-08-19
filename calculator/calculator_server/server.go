package main

import (
	"calculator-grpc/calculatorpb"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("funcion sum invocada con req %v\n", req)
	res := req.GetFirstNumber() + req.GetSecondNumber()
	result := &calculatorpb.SumResponse{
		SumResult: res,
	}
	return result, nil
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
