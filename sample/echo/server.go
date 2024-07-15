package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/context"

	echo "github.com/appnet-org/golib/sample/echo-pb"
	"google.golang.org/grpc"
)

type server struct {
	echo.UnimplementedEchoServiceServer
}

func (s *server) Echo(ctx context.Context, x *echo.Msg) (*echo.Msg, error) {
	log.Printf("Server got: [%s]", x.GetBody())

	// Check if the message contains "sleep"
	if x.GetBody() == "sleep" {
		log.Printf("Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)
	}

	hostname, _ := os.Hostname()
	appendedBody := fmt.Sprintf("You've hit %s\n", hostname)
	msg := &echo.Msg{
		Body: appendedBody,
	}
	return msg, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	fmt.Printf("Starting server pod at port 9000\n")

	echo.RegisterEchoServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
