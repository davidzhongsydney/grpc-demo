package main

import (
	"context"
	"log"
	"net"
	"net/http"

	pb "gRPC-demo/model"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// implement the RouteGuideServer interface
type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	point.Latitude = 100
	point.Longitude = 200
	return &pb.Feature{Name: "test", Location: point}, nil
}

func main() {

	msg := make(chan string)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
		}
	}()

	gwmux := runtime.NewServeMux()

	err = pb.RegisterRouteGuideHandlerServer(context.Background(), gwmux, &routeGuideServer{})

	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	go func() {
		log.Fatalln(http.ListenAndServe(":8080", gwmux))
		// log.Fatalln(http.ListenAndServe(":9000", gwmux))
	}()

	<-msg
}
