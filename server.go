package main

import (
	context "context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/shirou/grpc-sample-server/pb"
)

// SampleServer implements `service Sample`.
type SampleServer struct {
}

// NewSampleServer returns a new SampleServer.
func NewSampleServer() *pb.SampleService {
	return &pb.SampleService{
		Echo: Echo,
		Time: Time,
	}
}

// Echo implements Echo method
func Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoReply, error) {
	return &pb.EchoReply{Text: in.GetText()}, nil
}

// Time implements Time method
func Time(in *empty.Empty, stream pb.Sample_TimeServer) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			now := &pb.TimeMessage{
				Timestamp: ptypes.TimestampNow(),
			}
			if err := stream.Send(now); err != nil {
				return err
			}
		}
	}
}
