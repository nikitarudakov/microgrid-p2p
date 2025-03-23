package simulation

import (
	"context"
	"github.com/nikitarudakov/microenergy/internal/pb"
	"time"
)

type Meter struct {
	id  string
	cli pb.MeteringClient
	ctx context.Context
}

func NewMeter(id string, cli pb.MeteringClient, ctx context.Context) (*Meter, error) {
	return &Meter{
		id:  id,
		cli: cli,
		ctx: ctx,
	}, nil
}

func (m *Meter) Run() error {
	t := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-m.ctx.Done():
			return m.ctx.Err()
		default:
		}

		// Wait for 1 min to elapse
		<-t.C

		if _, err := m.cli.UploadMeteringReading(context.Background(), &pb.Reading{
			MeterId: m.id,
			Value:   122.3, // kW in a span of 1 min
		}); err != nil {
			return err
		}
	}
}
