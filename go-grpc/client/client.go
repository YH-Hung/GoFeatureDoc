package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "go-grpc/routeguide"
)

var serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")

func main() {
	flag.Parse()

	// Set up a connection to the gRPC server.
	conn, err := grpc.NewClient("dns:///"+*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// Close the connection after main function exits.
	defer conn.Close()

	// Create a new RouteGuide stub.
	client := pb.NewRouteGuideClient(conn)

	point := &pb.Point{Latitude: 409146138, Longitude: -746188906}
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)

	// Call GetFeature method on the client.
	// context.TODO() is a legitimate placeholder option when need to pass context.
	feature, err := client.GetFeature(context.TODO(), point)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %v", err)
	}
	log.Println(feature)

	///////////////////////////////////////////////////////////////////////////
	// Server-to-Client Streaming RPC
	//
	// Call ListFeatures method on the client.
	///////////////////////////////////////////////////////////////////////////
	rect := &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	}
	log.Printf("Looking for features within %v", rect)
	stream, err := client.ListFeatures(context.Background(), rect)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}
	for {
		// For server-to-client streaming RPCs, you call stream.Recv() until it
		// returns io.EOF.
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures failed: %v", err)
		}
		log.Printf("Feature: name: %q, point:(%v, %v)", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}

	///////////////////////////////////////////////////////////////////////////
	// Client-to-Server Streaming RPC
	//
	// Call RecordRoute method on the client.
	///////////////////////////////////////////////////////////////////////////
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))
	// Call RecordRoute method on the client.
	c2sStream, err := client.RecordRoute(context.Background())
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}
	// Stream each point to the server.
	for _, point := range points {
		if err := c2sStream.Send(point); err != nil {
			log.Fatalf("client.RecordRoute: stream.Send(%v) failed: %v", point, err)
		}
	}
	// Close the stream and receive the RouteSummary from the server.
	reply, err := c2sStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}
	log.Printf("Route summary: %v", reply)

	///////////////////////////////////////////////////////////////////////////
	// Bidirectional Streaming RPC
	//
	// Call RouteChat method on the client.
	///////////////////////////////////////////////////////////////////////////
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
	// Call RouteChat method on the client.
	biDiStream, err := client.RouteChat(context.Background())
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	// this channal is used to wait for the receive goroutine to finish.
	recvDoneCh := make(chan struct{})
	// receive goroutine.
	go func() {
		for {
			in, err := biDiStream.Recv()
			if err == io.EOF {
				// read done.
				close(recvDoneCh)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	// send messages simultaneously.
	for _, note := range notes {
		if err := biDiStream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}
	biDiStream.CloseSend()
	// wait for the receive goroutine to finish.
	<-recvDoneCh
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}
