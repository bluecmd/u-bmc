// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/u-root/u-bmc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type mgmtServer struct {

}

func (m *mgmtServer) PressButton(ctx context.Context, r *pb.ButtonPressRequest) (*pb.ButtonPressResponse, error) {
	err := PressButton(r.Button, r.DurationMs)
	if err != nil {
		return nil, err
	}
	return &pb.ButtonPressResponse{}, nil
}

func (m *mgmtServer) GetFans(ctx context.Context, _ *pb.GetFansRequest) (*pb.GetFansResponse, error) {
	// TODO(bluecmd): Implement
	return nil, fmt.Errorf("Unimplemented")
}

func (m *mgmtServer) SetFan(ctx context.Context, _ *pb.SetFanRequest) (*pb.SetFanResponse, error) {
	// TODO(bluecmd): Implement
	return nil, fmt.Errorf("Unimplemented")
}

func streamIn(stream pb.ManagementService_StreamConsoleServer) {
	for {
		buf := <-uartIn
		stream.Send(&pb.ConsoleData{Data: buf})
	}
}

func (m *mgmtServer) StreamConsole(stream pb.ManagementService_StreamConsoleServer) error {
	go streamIn(stream)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		uartOut <- in.Data
	}
}

func startGrpc() {
	// TODO(bluecmd): Since we have no RNG, no configuration, etc
	// only use http for now
	l, err := net.Listen("tcp", "[::]:80")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	s := mgmtServer{}

	g := grpc.NewServer()
	pb.RegisterManagementServiceServer(g, &s)
	reflection.Register(g)
	go g.Serve(l)
}
