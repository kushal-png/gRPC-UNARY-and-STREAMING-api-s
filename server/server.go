package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "grpc/proto"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type Server struct {
	pb.UnimplementedGreetServiceServer
}

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start the server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &Server{})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
	log.Printf("Server started at %v:", listener.Addr())
}

func (s *Server) SayHello(ctx context.Context, req *pb.NoParams) (*pb.HelloResponse, error) {
	res := &pb.HelloResponse{Message: "Hello from Server"}
	return res, nil
}

func (s *Server) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("got request with names: %v", req.Names)
	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}



func (s *Server) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var message []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessageList{Messages: message})
		}
		if err != nil {
			return err
		}
		log.Printf("Got Request with name: %v", req.Name)
		message= append(message, "Hello "+req.Name)
	}
}

func (s *Server) SayHelloBidirectionStreaming(stream pb.GreetService_SayHelloBidirectionStreamingServer) error {
	
	for{
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got Request with name: %v", req.Name)

		res:= &pb.HelloResponse{
			Message: "Hello "+ req.Name ,
		}
        
		if err:= stream.Send(res); err!=nil{
			return err
		}
	}
}
