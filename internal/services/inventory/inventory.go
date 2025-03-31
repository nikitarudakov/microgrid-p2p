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
	ImportMeterID     uuid.UUID `json:"import_meter_id"`
	ExportMeterID     uuid.UUID `json:"export_meter_id"`
	ProviderID        uuid.UUID `json:"provider_id"`
	VoltageLevel      float32   `json:"voltage_level"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	MaxCapacity       float32   `json:"max_capacity"`
	MinRuntimeMinutes int32     `json:"min_runtime_minutes"`
	MaxRuntimeMinutes int32     `json:"max_runtime_minutes"`
}

type Service struct {
	db *gorm.DB
	pb.UnimplementedInventoryManagementServer
}

func (s *Service) RegisterAsset(_ context.Context, in *pb.Asset) (*pb.Asset, error) {
	asset := pb.FromProto(in, &Asset{})

	if err := s.db.Create(asset).Error; err != nil {
		return nil, err
	}

	return pb.ToProto(asset, &pb.Asset{}), nil
}
