package main

import (
	"fmt"
	v1 "github.com/nikitarudakov/microenergy/internal/gen/inventory/v1"
	"github.com/nikitarudakov/microenergy/internal/services/inventory"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// Set up logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Build the Host Connection String
	host := fmt.Sprintf(":%s", "5000")

	// Create the TCP Listener
	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.Fatalln("unable to start TCP server:", err.Error())
	}

	// Set the gRPC Options
	var opts []grpc.ServerOption

	// Create the new gRPC Server
	grpcServer := grpc.NewServer(opts...)

	// Register the gRPC Server
	v1.RegisterInventoryManagementServer(grpcServer, &inventory.Server{})

	// Serve the Servants, oh no
	logger.Println("Starting Company Service")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalln("unable to start gRPC server:", err.Error())
	}
}
