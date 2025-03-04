package inventory

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Server struct {
	pb.InventoryManagementServer
}

func (s *Server) GetOwnerEnergyResourceList(_ context.Context, ownerId *wrapperspb.StringValue) (*pb.EnergyResourceList, error) {
	return &pb.EnergyResourceList{
		EnergyResources: []*pb.EnergyResource{{
			Id:       uuid.New().String(),
			OwnerId:  ownerId.String(),
			Capacity: 123.24,
		}},
	}, nil
}
