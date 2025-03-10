package pb

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func Connect(service string) (*grpc.ClientConn, error) {
	host := os.Getenv(fmt.Sprintf("%s_HOST", service))
	port := os.Getenv(fmt.Sprintf("%s_PORT", service))

	connString := fmt.Sprintf("%s:%s", host, port)

	return grpc.NewClient(
		connString,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
}
