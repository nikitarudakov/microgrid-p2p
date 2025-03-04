package inventory

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sync"
)

type Server struct {
	energyResources []*pb.EnergyResource
	mu              sync.Mutex
	pb.InventoryManagementServer
}

func (s *Server) GetOwnerEnergyResourceList(_ context.Context, ownerName *wrapperspb.StringValue) (*pb.EnergyResourceList, error) {
	var ownersEnergyResources []*pb.EnergyResource
	for _, er := range s.energyResources {
		if er.OwnerName == ownerName.Value {
			ownersEnergyResources = append(ownersEnergyResources, er)
		}
	}

	return &pb.EnergyResourceList{EnergyResources: ownersEnergyResources}, nil
}

func (s *Server) RegisterEnergyResource(_ context.Context, in *pb.RegisterEnergyResourceInput) (*pb.EnergyResource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	energyResource := &pb.EnergyResource{
		Id:        uuid.New().String(),
		OwnerName: in.OwnerName,
		Capacity:  in.Capacity,
	}

	s.energyResources = append(s.energyResources, energyResource)

	return energyResource, nil
}
