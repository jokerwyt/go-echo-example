package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	echo "github.com/appnet-org/golib/sample/echo-stream/pb"
	"google.golang.org/grpc"
)

var (
	grpcClient echo.EchoServiceClient
	stream     echo.EchoService_EchoClient
)

func initGRPC() {
	conn, err := grpc.Dial("server-stream.default.svc.cluster.local:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	grpcClient = echo.NewEchoServiceClient(conn)
	stream, err = grpcClient.Echo(context.Background())
	if err != nil {
		log.Fatalf("could not establish stream: %v", err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive a message: %v", err)
			}
			log.Printf("Got message: %s", in.Body)
		}
	}()
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody := r.URL.Query().Get("key")
	if requestBody == "" {
		http.Error(w, "key parameter is required", http.StatusBadRequest)
		return
	}

	for i := 0; i < 2; i++ {
		msg := fmt.Sprintf("%s %d", requestBody, i+1)
		if err := stream.Send(&echo.Msg{Body: msg}); err != nil {
			log.Fatalf("Failed to send a message: %v", err)
		}
		log.Printf("Sent message: %s", msg)
		time.Sleep(1 * time.Second) // Simulate some delay
	}

	fmt.Fprintf(w, "Sent message: %s", requestBody)
}

func main() {
	initGRPC()

	http.HandleFunc("/", handleHTTP)
	log.Println("HTTP server listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}

	// Clean up the stream after the HTTP server is done
	stream.CloseSend()
	time.Sleep(1 * time.Second) // Allow some time to finish receiving messages
}
