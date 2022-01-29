package main

import (
	"context"
	"fmt"
	"log"

	"github.com/8thgencore/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC... ")
	req := &calculatorpb.SumRequest{
		FirstNumber:  15,
		SecondNumber: 18,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.SumResult)
}
