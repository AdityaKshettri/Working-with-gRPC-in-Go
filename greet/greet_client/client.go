package main

import (
	"fmt"
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
}