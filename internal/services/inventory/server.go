package inventory

import (
	"context"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sync"
)

type Server struct {
	energyResources []*pb.EnergyResource
	mu              sync.Mutex
	pb.InventoryManagementServer
}

func (s *Server) FetchAllEnergyResources(_ context.Context, _ *emptypb.Empty) (*pb.EnergyResourceList, error) {
	return &pb.EnergyResourceList{EnergyResources: s.energyResources}, nil
}

func (s *Server) FetchProducerEnergyResources(_ context.Context, producerId *wrapperspb.StringValue) (*pb.EnergyResourceList, error) {
	var ownersEnergyResources []*pb.EnergyResource
	for _, er := range s.energyResources {
		if er.ProducerId == producerId.Value {
			ownersEnergyResources = append(ownersEnergyResources, er)
		}
	}

	return &pb.EnergyResourceList{EnergyResources: ownersEnergyResources}, nil
}

func (s *Server) RegisterEnergyResource(_ context.Context, in *pb.RegisterEnergyResourceInput) (*pb.EnergyResource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	energyResource := &pb.EnergyResource{
		Id:         in.Id,
		Name:       in.Name,
		ProducerId: in.ProducerId,
		Capacity:   in.Capacity,
		Price:      in.Price,
	}

	s.energyResources = append(s.energyResources, energyResource)

	return energyResource, nil
}
