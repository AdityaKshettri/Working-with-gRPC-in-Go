package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/AdityaKshettri/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello! I am a client!")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created Client : %f", c)

	//doUnary(c)
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Kshettri",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC : %v", err)
	}
	log.Printf("Response from Greet : %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Kshettri",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC : %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// We have reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream : %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}
