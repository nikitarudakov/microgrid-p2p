package bidding

import (
	"cmp"
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/services/competition"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"slices"
)

type Bid struct {
	ServiceWindowID   uuid.UUID `json:"service_window_id"`
	AssetID           uuid.UUID `json:"asset_id"`
	Capacity          float64   `json:"capacity"`
	AvailabilityPrice float64   `json:"availability_price"`
	UtilizationPrice  float64   `json:"utilization_price"`
	ServiceFee        float64   `json:"service_fee"`
	RuntimeMinutes    int32     `json:"runtime_minutes"`
}

type Service struct {
	db *mongo.Database
}

type Winner struct {
	Bid        *Bid
	Window     *competition.Window
	TotalPrice float64 // For tracking cost
}

func (s *Service) SelectWinners(comp *competition.Competition, allBids []*Bid) []*Winner {
	var winners []*Winner
	usedAssets := map[uuid.UUID]bool{}

	for _, window := range comp.Windows {
		windowBids := filterBidsForWindow(allBids, window)

		// Sort bids by total price (lowest first)
		slices.SortFunc(windowBids, func(a, b *Bid) int {
			costA := totalBidCost(a)
			costB := totalBidCost(b)
			return cmp.Compare(costA, costB)
		})

		var totalCapacity float64
		var selected []*Winner

		for _, bid := range windowBids {
			if usedAssets[bid.AssetID] {
				continue
			}

			if totalCapacity >= float64(window.Capacity) {
				break
			}

			cost := totalBidCost(bid)
			if comp.MaxBudget > 0 && (sumWinnerPrices(selected)+cost) > float64(comp.MaxBudget) {
				continue // Over budget
			}

			selected = append(selected, &Winner{
				Bid:        bid,
				Window:     window,
				TotalPrice: cost,
			})
			usedAssets[bid.AssetID] = true
			totalCapacity += bid.Capacity
		}

		winners = append(winners, selected...)
	}

	return winners
}
