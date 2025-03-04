package resolver

import (
	inventoryService "github.com/nikitarudakov/microenergy/internal/pb"
	"github.com/sirupsen/logrus"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	inventoryManagementService inventoryService.InventoryManagementClient
	logger                     *logrus.Logger
}

func NewResolver(
	inventory inventoryService.InventoryManagementClient,
	logger *logrus.Logger,
) *Resolver {
	return &Resolver{
		inventoryManagementService: inventory,
		logger:                     logger,
	}
}
