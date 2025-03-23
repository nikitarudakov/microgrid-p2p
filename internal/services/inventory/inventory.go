package inventory

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"gorm.io/gorm"
)

type Asset struct {
	ID                uuid.UUID `json:"id"`
	Ref               string    `json:"ref"`
	Name              string    `json:"name"`
	ImportMeterID     string    `json:"import_meter_id"`
	ExportMeterID     string    `json:"export_meter_id"`
	ProviderID        string    `json:"provider_id"`
	MaxCapacity       float32   `json:"max_capacity"`
	VoltageLevel      float32   `json:"voltage_level"`
	MaxRuntimeMinutes int32     `json:"max_runtime_minutes"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	Services          []string  `json:"services"`
}

type Server struct {
	db *gorm.DB
	pb.UnimplementedInventoryManagementServer
}

func (s *Server) RegisterAsset(_ context.Context, in *pb.Asset) (*pb.Asset, error) {
	asset := pb.FromProto(in, &Asset{})

	if err := s.db.Create(asset).Error; err != nil {
		return nil, err
	}

	return pb.ToProto(asset, &pb.Asset{}), nil
}
