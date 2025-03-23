package metering

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/onchain"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math"
	"time"
)

// Reading represents a NoSQL structure.
type Reading struct {
	RecordedAt   time.Time `json:"recorded_at" bson:"recorded_at"`
	ObligationID string    `json:"dispatch_id" bson:"dispatch_id"`
	MeterID      string    `json:"meter_id" bson:"meter_id"`
	Interval     float64   `json:"interval" bson:"interval"`
	Value        float64   `json:"value" bson:"value"`
	Baseline     float64   `json:"baseline" bson:"baseline"`
	Direction    string    `json:"direction" bson:"direction"`
}

type Server struct {
	db *mongo.Database
	pb.UnimplementedMeteringServer
}

func (s *Server) UploadMeteringReading(_ context.Context, in *pb.Reading) (_ *emptypb.Empty, err error) {
	coll := s.db.Collection("readings")
	reading := pb.FromProto(in, &Reading{})

	_, err = coll.InsertOne(context.Background(), reading)

	return
}

func (s *Server) RecordDispatch(_ context.Context, in *pb.RecordDispatchRequest) (*pb.RecordDispatchResponse, error) {
	coll := s.db.Collection("readings")

	filter := bson.D{
		{"dispatch_id", in.DispatchId},
		{"recorded_at", bson.D{
			{"$gte", in.StartTime.AsTime()},
			{"$lte", in.EndTime.AsTime()},
		}},
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var readings []Reading
	if err = cursor.All(context.Background(), &readings); err != nil {
		return nil, err
	}

	if len(readings) < 1 {
		return nil, errors.New("empty dispatch interval readings slice")
	}

	dispatch := &onchain.Dispatch{
		ID:           uuid.New().String(),
		DocType:      "dispatch",
		ObligationID: readings[0].ObligationID,
		TimeWindow: onchain.TimeWindow{
			StartTime: in.StartTime.AsTime().Format(time.RFC3339),
			EndTime:   in.EndTime.AsTime().Format(time.RFC3339),
		},
		Direction: readings[0].Direction,
		Capacity:  0,
	}

	for _, reading := range readings {
		if reading.Direction != dispatch.Direction {
			log.Printf("invalid interval reading direction - curr: %s != dispatched: %s\n", reading.Direction, dispatch.Direction)
			continue
		}

		switch dispatch.Direction {
		case "import":
			dispatch.Capacity += math.Abs(reading.Baseline-reading.Value) * reading.Interval // (120 - 60) * 0.25 = 12 kWh (delivered)
		case "export":
			dispatch.Capacity += reading.Value * reading.Interval
		}
	}

	// TODO: Register on blockchain with blockchain client

	return &pb.RecordDispatchResponse{}, nil
}
