package inventory

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/gen/inventory/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Server struct {
	v1.InventoryManagementServer
}

func (s *Server) GetOwnerEnergyResourceList(_ context.Context, ownerId *wrapperspb.StringValue) (*v1.EnergyResourceList, error) {
	return &v1.EnergyResourceList{
		EnergyResources: []*v1.EnergyResource{{
			Id:       uuid.New().String(),
			OwnerId:  ownerId.String(),
			Capacity: 123.24,
		}},
	}, nil
}
