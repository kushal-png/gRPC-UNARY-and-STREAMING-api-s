package main

import (
	"context"
	"fmt"
	pb "grpc/proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func CallSayHello(client pb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &pb.NoParams{})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Printf("%s", res.Message)
}

func CallSayHelloServerStreaming(client pb.GreetServiceClient, names *pb.NamesList) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fmt.Println("streaming started")

	stream, err := client.SayHelloServerStreaming(ctx, names)
	if err != nil {
		log.Fatalf("could not send names %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("Ended")
		}
		if err != nil {
			log.Fatalf("error while streaming: %v", err)
		}
		log.Println(message)
	}

}

func CallSayHelloClientStreaming(client pb.GreetServiceClient, names *pb.NamesList) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Client streaming started")

	stream, err := client.SayHelloClientStreaming(ctx)
	if err != nil {
		log.Fatalf("Error while streaming: %v", err)
	}

	for _, name := range names.Names {
		req := &pb.HelloRequest{Name: name}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending: %v", err)
		}
		log.Printf("Sent the request with name: %s", name)
		time.Sleep(2 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("CLient srtreaming error")
	}
	log.Printf("Client Streaming Finished \n")
	log.Printf("%v", res.Messages)

}

func CallSayHelloBidirectionalStreaming(client pb.GreetServiceClient, names *pb.NamesList) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Bidirectional streaming started")
	stream, err := client.SayHelloBidirectionStreaming(ctx)
	if err != nil {
		log.Fatalf("Cannot Send names")
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while streaming: %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{Name: name}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional Streaming Finished")
}

func main() {
	connection, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer connection.Close()

	client := pb.NewGreetServiceClient(connection)
	names := &pb.NamesList{Names: []string{"Akhil", "Kushal", "Kunal", "keshav", "naman"}}
	CallSayHelloBidirectionalStreaming(client, names)
}
