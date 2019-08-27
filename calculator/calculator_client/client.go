package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/crojas/go-grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDirectionalStreaming(c)
	doErrorUnary(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error llamando a funcion ComputeAverage %v", err)
	}
	for i := 0; i < 5; i++ {
		randomNumber := (rand.Float64() * float64(i)) + float64(i)
		req := &calculatorpb.ComputeAverageRequest{
			Number: randomNumber,
		}
		fmt.Printf("Enviando %v\n", randomNumber)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error recibiendo mensage del server %v", err)
	}
	fmt.Printf("El promedio es: %v", res.GetAverage())
}

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error llamando a funcion FindMaximum %v", err)
	}
	waitc := make(chan struct{})

	// Enviando numeros en una goroutine.
	go func() {
		for i := 0; i < 10; i++ {
			number := rand.Int63n(100)
			fmt.Printf("Enviando %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Recibiendo el numero maximo.
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error recibiendo numero maximo %v\n", err)
			}
			fmt.Printf("Maximo de los enviados: %v\n", res.GetMax())
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: -1,
	})
	if err != nil {
		respError, ok := status.FromError(err)
		if ok {
			// Error de GRPC
			fmt.Println(respError.Message())
			fmt.Println(respError.Code())
			if respError.Code() == codes.InvalidArgument {
				fmt.Println("Enviamos un numero negativo!")
				return
			}
		} else {
			log.Fatalf("Error llamando a SquareRoot %v", err)
			return
		}
	}
	fmt.Printf("Respuesta %v", res.GetNumberRoot())
}
