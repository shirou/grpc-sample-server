package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shirou/grpc-sample-server/pb"

	"google.golang.org/grpc"
)

var (
	serverAddr string
)

func main() {
	flag.StringVar(&serverAddr, "addr", "localhost:8888", "The server address in the format of host:port")
	// convert Environment Variables to flags
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()

	start()
}

func start() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewSampleClient(conn)

	echo(client, "ECHO ECHO")

	timeStream(client)

}

func echo(client pb.SampleClient, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	reply, err := client.Echo(ctx, &pb.EchoRequest{Text: text})
	if err != nil {
		log.Fatalf("%v.Echo(_) = _, %v: ", client, err)
	}

	fmt.Println(reply)
}

func timeStream(client pb.SampleClient) {
	// Will be DeadlineExceeded after 100 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	stream, err := client.Time(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("%v.Time(_) = _, %v: ", client, err)
	}
	for {
		serverTime, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Time(_) = _, %v", client, err)
		}
		st, err := ptypes.Timestamp(serverTime.GetTimestamp())
		if err != nil {
			log.Fatalf("%v.TimeStamp(_) = _, %v", client, err)
		}
		now := time.Now()
		diff := now.Sub(st)

		fmt.Printf("server: %s, client: %s, diff: %s\n", st, now, diff)
	}
}
