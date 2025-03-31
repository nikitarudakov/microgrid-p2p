package onchain

import (
	"fmt"
	"time"
)

type OutOfRangeTimeWindowError struct {
	startC string
	endC   string
	startO string
	endO   string
}

func (e *OutOfRangeTimeWindowError) Error() string {
	return fmt.Sprintf("violated time window - contracted: %s - %s does not contain obliged: %s - %s", e.startC, e.endC, e.startO, e.endO)
}

type RequestedCapacityExceededError struct {
	contracted float64
	requested  float64
}

func (e *RequestedCapacityExceededError) Error() string {
	return fmt.Sprintf("exceeded capacity - contracted: %.2f < requsted: %2.f", e.contracted, e.requested)
}

type InsufficientCapacityError struct {
	dispatched float64
	needed     float64
}

func (e *InsufficientCapacityError) Error() string {
	return fmt.Sprintf("insufficient capacity - dispatched: %.2f < needed: %2.f", e.dispatched, e.needed)
}

type IncorrectDirectionError struct {
	dispatched string
	needed     string
}

func (e *IncorrectDirectionError) Error() string {
	return fmt.Sprintf("incorrect direction - dispatched %s != obliged: %s", e.dispatched, e.needed)
}

type ViolatedTimeWindowError struct {
	startD string
	endD   string
	startO string
	endO   string
}

func (e *ViolatedTimeWindowError) Error() string {
	return fmt.Sprintf("violated time window - dispatch: %s - %s != obliged: %s - %s", e.startD, e.endD, e.startO, e.endO)
}

// TimeWindow defines time brackets with StartTime and EndTime
type TimeWindow struct {
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Baseline  float64 `json:"baseline"`
}

func (tw TimeWindow) Hours() (float64, error) {
	layout := time.RFC3339

	startTime, err := time.Parse(layout, tw.StartTime)
	if err != nil {
		return 0, err
	}

	endTime, err := time.Parse(layout, tw.EndTime)
	if err != nil {
		return 0, err
	}

	// Calculate the duration and convert to hours
	duration := endTime.Sub(startTime)

	return duration.Hours(), nil
}

func (tw TimeWindow) Parse() (time.Time, time.Time, error) {
	layout := time.RFC3339

	startTime, err := time.Parse(layout, tw.StartTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time: %w", err)
	}

	endTime, err := time.Parse(layout, tw.EndTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %w", err)
	}

	return startTime, endTime, nil
}

func (tw TimeWindow) Contains(tw2 TimeWindow) bool {
	start1, end1, err := tw.Parse()
	if err != nil {
		return false
	}

	start2, end2, err := tw2.Parse()
	if err != nil {
		return false
	}

	return (start1.Before(start2) || start1.Equal(start2)) && (end1.After(end2) || end1.Equal(end2))
}

// Contract is based on a competition details and a winning bid.
type Contract struct {
	// Unique contract identifier
	ID      string `json:"id"`
	DocType string `json:"doc_type"`
	// Parties data
	BuyerID       string `json:"buyer_id"`
	SellerID      string `json:"seller_id"`
	CompetitionID string `json:"competition_id"`
	// Contract's terms
	StartDate  string     `json:"start_dates"`
	EndDate    string     `json:"end_date"`
	Penalty    float64    `json:"penalty"`
	TimeWindow TimeWindow `json:"time_window"`
	// Bid Data
	Capacity          float64 `json:"capacity"`
	AvailabilityPrice float64 `json:"availability_price"`
	UtilizationPrice  float64 `json:"utilization_price"`
	ServiceFee        float64 `json:"service_fee"`
	// Might be a single asset or a group of assets
	Assets []string `json:"assets"`
}

// Obligation takes hold after FSP accepts dispatch request
// - Status: if unfulfilled an FSP is a subject to penalty, in any case settle for a provided flexibility.
type Obligation struct {
	ID         string     `json:"id"`
	DocType    string     `json:"doc_type"`
	ContractID string     `json:"contract_id"`
	TimeWindow TimeWindow `json:"time_window"`
	Direction  string     `json:"direction"` // import or export
	Status     string     `json:"status"`    // fulfilled, unfulfilled, stop_requested
	Capacity   float64    `json:"capacity"`
}

func (o Obligation) Validate(c *Contract) error {
	if !c.TimeWindow.Contains(o.TimeWindow) {
		return &OutOfRangeTimeWindowError{
			startC: c.TimeWindow.StartTime,
			endC:   c.TimeWindow.EndTime,
			startO: o.TimeWindow.StartTime,
			endO:   o.TimeWindow.EndTime,
		}
	}

	if c.Capacity > o.Capacity {
		return &RequestedCapacityExceededError{
			contracted: c.Capacity,
			requested:  o.Capacity,
		}
	}

	return nil
}

// Dispatch is recorded at the end of obligation time brackets.
// It is further validated against an Obligation.
type Dispatch struct {
	ID           string     `json:"id"`
	DocType      string     `json:"doc_type"`
	ObligationID string     `json:"obligation_id"`
	TimeWindow   TimeWindow `json:"time_window"`
	Direction    string     `json:"direction"`
	Capacity     float64    `json:"capacity"`
}

func (d *Dispatch) CalculatePayableAmount(c *Contract) float64 {
	if c.ServiceFee != 0 {
		return c.ServiceFee
	}

	var totalAmount float64
	hours, _ := d.TimeWindow.Hours()

	// Add availability amount
	totalAmount += c.Capacity * c.AvailabilityPrice * hours

	// Add dispatch amount
	totalAmount += d.Capacity * c.UtilizationPrice

	return totalAmount
}

func (d *Dispatch) Validate(o *Obligation) error {
	if d.Capacity < o.Capacity {
		return &InsufficientCapacityError{dispatched: d.Capacity, needed: o.Capacity}
	}

	if o.Direction != d.Direction {
		return &IncorrectDirectionError{dispatched: d.Direction, needed: o.Direction}
	}

	startD, endD, err := d.TimeWindow.Parse()
	if err != nil {
		return err
	}

	startO, endO, err := o.TimeWindow.Parse()
	if err != nil {
		return err
	}

	// If was started after obliged or ended before obliged
	if startD.After(startO) || endD.Before(endO) {
		return &ViolatedTimeWindowError{
			startD: d.TimeWindow.StartTime, endD: d.TimeWindow.EndTime,
			startO: d.TimeWindow.StartTime, endO: d.TimeWindow.StartTime,
		}
	}

	return nil
}

type Settlement struct {
	ID         string  `json:"ID"`
	DispatchID string  `json:"dispatch_id"`
	Type       string  `json:"type"` // payable or penalised
	Fiat       float64 `json:"fiat"`
	SettledAt  string  `json:"settled_at"`
}
