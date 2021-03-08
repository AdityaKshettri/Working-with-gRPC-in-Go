package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/AdityaKshettri/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + "number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with client stream request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We have finished reading the client stream
			res := &greetpb.LongGreetResponse{
				Result: result,
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked with BiDi stream request")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		res := &greetpb.GreetEveryoneResponse{
			Result: result,
		}
		sendErr := stream.Send(res)
		if sendErr != nil {
			log.Fatalf("Error while sending data to client : %v", sendErr)
			return sendErr
		}
	}

}

func main() {
	fmt.Println("Hello World!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server : %v", err)
	}
}

func (*server) GreetWithDealine(ctx context.Context, req *greetpb.GreetWithDealineRequest) (*greetpb.GreetWithDealineResponse, error) {
	fmt.Printf("GreetWithDealine function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Client cancelled the request!")
			return nil, status.Error(codes.Canceled, "Client cancelled the request!")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDealineResponse{
		Result: result,
	}
	return res, nil
}
