package contracts

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"log"
	"time"
)

var (
	ErrRequestCapacityLimitExceeded = errors.New("requested capacity exceeds an obligated capacity")
	ErrRequestNotInServiceWindow    = errors.New("request is out of service window time range")
)

// Voltage represent voltage requirements for a competition
type Voltage struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

// TimeWindow is used to describe Dispatch requested time window
// Start - format:
type TimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Runtime struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (tw TimeWindow) parseForToday() (time.Time, time.Time, error) {
	start, err := parseServiceWindowEdge(tw.Start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := parseServiceWindowEdge(tw.End)
	if err != nil {
		return start, time.Time{}, err
	}
	return start, end, err
}

func (tw TimeWindow) Includes(tw2 TimeWindow) bool {
	start1, end1, err := tw.parseForToday()
	if err != nil {
		log.Printf("error parsing time window: %s - %s", tw.Start, tw.End)
		return false
	}

	start2, end2, err := tw2.parseForToday()
	if err != nil {
		log.Printf("error parsing time window: %s - %s", tw2.Start, tw2.End)
		return false
	}

	// If overlaps return false
	if start2.Before(start1) || end2.After(end1) {
		return false
	}

	// Completely includes a time interval
	return true
}

// Dispatch represents recorded dispatch of FSP
type Dispatch struct {
	ObligationID string  `json:"obligation_id"`
	RequestID    string  `json:"request_id"`
	Flexibility  float32 `json:"capacity"`
	Runtime      Runtime `json:"runtime"`
}

// PricingType is used to describe a pricing used in the FSP's bid
type PricingType struct {
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

// Request represents a dispatch requested by a consumer
// It is a key structure for verifying agreement fulfillment.
type Request struct {
	ID           string     `json:"id"`
	ObligationID string     `json:"obligation_id"`
	Capacity     float32    `json:"capacity"`
	Status       string     `json:"status"`
	TimeWindow   TimeWindow `json:"time_window"`
}

// Obligation represents an agreed deal between FSP and SO.
// It stores all the data needed for validating FSPs dispatches and their requests from SOs.
type Obligation struct {
	ID            string     `json:"id"`
	Capacity      float32    `json:"capacity"`
	ServiceWindow TimeWindow `json:"service_window"`
	Frequency     string     `json:"frequency"`
	MeterID       string     `json:"meter_id"`
}

func (a Obligation) validateRequest(req Request) error {
	if !a.ServiceWindow.Includes(req.TimeWindow) {
		return ErrRequestNotInServiceWindow
	}

	// Check if capacity does not exceed agreed one
	if req.Capacity > a.Capacity {
		return ErrRequestCapacityLimitExceeded
	}

	return nil
}

// Fulfillment works on Hyperledger smart contract technology.
// It automatically handles verification of dispatch and billing.
type Fulfillment struct {
	contractapi.Contract
}

// RegisterObligation stores agreement which contains information about the competition
// and its bids. This method should be used after competition auction has finished and
// a set of bids that fulfill competition requirement were accepted.
func (f *Fulfillment) RegisterObligation(
	ctx contractapi.TransactionContextInterface,
	obligation *Obligation,
) error {
	stub := ctx.GetStub()

	data, err := json.Marshal(obligation)
	if err != nil {
		return err
	}

	return stub.PutState(obligation.ID, data)
}

// RequestDispatch validates a dispatch request of System Operator and updates
// agreement's requests slice. The requests slice can be used for audit purposes.
func (f *Fulfillment) RequestDispatch(ctx contractapi.TransactionContextInterface, req Request) error {
	stub := ctx.GetStub()

	data, err := stub.GetState(req.ObligationID)
	if err != nil {
		return err
	}

	obligation := &Obligation{}
	if err := json.Unmarshal(data, obligation); err != nil {
		return nil
	}

	if err := obligation.validateRequest(req); err != nil {
		return err
	}

	// Emit event that is used to track metering data
	if err := stub.SetEvent("DispatchRequested", []byte(obligation.MeterID)); err != nil {
		return err
	}

	data, err = json.Marshal(obligation)
	if err != nil {
		return err
	}

	return stub.PutState(obligation.ID, data)
}
