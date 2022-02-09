package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/8thgencore/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC... ")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lil",
			LastName:  "Wayne",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC... ")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lil",
			LastName:  "Wayne",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC... ")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	// requests := []*greetpb.LongGreetRequest{
	// 	&greetpb.LongGreetRequest{
	// 		Greeting: &greetpb.Greeting{
	// 			FirstName: "Stephane",
	// 		},
	// 	},
	// 	&greetpb.LongGreetRequest{
	// 		Greeting: &greetpb.Greeting{
	// 			FirstName: "John",
	// 		},
	// 	},
	// 	&greetpb.LongGreetRequest{
	// 		Greeting: &greetpb.Greeting{
	// 			FirstName: "Lucy",
	// 		},
	// 	},
	// 	&greetpb.LongGreetRequest{
	// 		Greeting: &greetpb.Greeting{
	// 			FirstName: "Mark",
	// 		},
	// 	},
	// 	&greetpb.LongGreetRequest{
	// 		Greeting: &greetpb.Greeting{
	// 			FirstName: "Piper",
	// 		},
	// 	},
	// }

	// for _, req := range requests {
	// 	log.Printf("Sending req: %v\n", req)
	// 	stream.Send(req)
	// 	time.Sleep(1000 * time.Millisecond)
	// }

	persons := []string{"Stephane", "John", "Lucy", "Mark", "Piper"}

	for _, person := range persons {
		log.Printf("Sending person: %v\n", person)
		stream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: person,
			},
		})
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)

	}
	fmt.Printf("LongGreet Response: %v\n", res)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC... ")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		persons := []string{"Stephane", "John", "Lucy", "Mark", "Piper"}

		for _, person := range persons {
			fmt.Printf("Sending person: %v\n", person)
			stream.Send(&greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: person,
				},
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			log.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
