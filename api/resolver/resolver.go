package resolver

import (
	"github.com/nikitarudakov/microenergy/internal/pb"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Services struct {
	inventoryManagementService pb.InventoryManagementClient
	userManagementService      pb.UserManagementClient
}

func (s *Services) Connect() error {
	grpcConn, err := pb.Connect("USER")
	if err != nil {
		return err
	}
	s.userManagementService = pb.NewUserManagementClient(grpcConn)

	grpcConn, err = pb.Connect("INVENTORY")
	if err != nil {
		return err
	}
	s.inventoryManagementService = pb.NewInventoryManagementClient(grpcConn)

	return nil
}

type Resolver struct {
	services *Services
	logger   *logrus.Logger
}

func New(logger *logrus.Logger) (*Resolver, error) {
	services := &Services{}
	if err := services.Connect(); err != nil {
		return nil, err
	}

	return &Resolver{
		services: services,
		logger:   logger,
	}, nil
}
