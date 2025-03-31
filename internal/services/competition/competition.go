package competition

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"gorm.io/gorm"
	"time"
)

type Window struct {
	ID                uuid.UUID    `json:"id"`
	CompetitionID     uuid.UUID    `json:"competition_id"`
	Competition       *Competition `json:"competition"`
	StartTime         time.Time    `json:"start_time"`
	EndTime           time.Time    `json:"end_time"`
	Weekdays          []string     `json:"weekdays"`
	Capacity          float32      `json:"capacity"`
	MinRuntimeMinutes int32        `json:"min_runtime_minutes"`
}

type Competition struct {
	ID          uuid.UUID `json:"id"`
	OrganizerID uuid.UUID `json:"organizer_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	MinVoltage  float32   `json:"min_voltage"`
	MaxVoltage  float32   `json:"max_voltage"`
	MaxBudget   float32   `json:"max_budget"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Radius      float64   `json:"radius"`
	Windows     []*Window `json:"windows"`
}

type Server struct {
	db *gorm.DB
	pb.UnimplementedCompetitionManagementServer
}

func (s *Server) RegisterCompetition(_ context.Context, in *pb.Competition) (*pb.Competition, error) {
	competition := pb.FromProto(in, &Competition{})

	if err := s.db.Create(competition).Error; err != nil {
		return nil, err
	}

	return pb.ToProto(competition, &pb.Competition{}), nil
}
