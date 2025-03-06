package resolver

import (
	"github.com/nikitarudakov/microenergy/internal/pb"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	inventoryManagementService pb.InventoryManagementClient
	userManagementService      pb.UserManagementClient
	logger                     *logrus.Logger
}

func NewResolver(
	inventory pb.InventoryManagementClient,
	user pb.UserManagementClient,
	logger *logrus.Logger,
) *Resolver {
	return &Resolver{
		inventoryManagementService: inventory,
		userManagementService:      user,
		logger:                     logger,
	}
}
