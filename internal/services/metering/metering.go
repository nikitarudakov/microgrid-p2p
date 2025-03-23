package metering

import (
	"context"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

// Reading represents a NoSQL structure
type Reading struct {
	RecordedAt time.Time `json:"recorded_at" bson:"recorded_at"`
	DispatchID string    `json:"dispatch_id" bson:"dispatch_id"`
	MeterID    string    `json:"meter_id" bson:"meter_id"`
	Interval   string    `json:"interval" bson:"interval"`
	Value      float64   `json:"value" bson:"value"`
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
