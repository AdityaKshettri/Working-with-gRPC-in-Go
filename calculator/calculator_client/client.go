package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AdityaKshettri/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator Client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	fmt.Printf("Created Client : %f\n", c)

	//calculateSum(c)
	//calculatePrimeNumberDecomposition(c)
	//computeAverage(c)
	//findMaximum(c)
	findSquareRoot(c)
}

func calculateSum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC : %v", err)
	}
	log.Printf("Response from Sum : %v", res.SumResult)
}

func calculatePrimeNumberDecomposition(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeNumberDecomposition RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC : %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Some Error happened : %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func computeAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage RPC...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}
	numbers := []int32{3, 4, 5, 6}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	fmt.Printf("The Average is: %v\n", res.GetAverage())
}

func findMaximum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a FindMaximum RPC...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		numbers := []int32{4, 7, 1, 5, 19, 34}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading Server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received new Maximum: %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func findSquareRoot(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot RPC...")
	//correct call
	findSquareRootForNumber(c, int32(10))
	//error call
	findSquareRootForNumber(c, int32(-10))
}

func findSquareRootForNumber(c calculatorpb.CalculatorServiceClient, n int32) {
	req := &calculatorpb.SquareRootRequest{
		Number: n,
	}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probabily sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big error while calling Square Root: %v\n", err)
			return
		}
	}
	fmt.Printf("Root of number %v: %v\n", n, res.GetNumberRoot())
}
