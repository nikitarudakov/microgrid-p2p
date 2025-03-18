package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/nikitarudakov/microenergy/internal/services/bidding"
	"github.com/nikitarudakov/microenergy/internal/services/inventory"
	"log"
)

var (
	ErrInvalidDispatchCapacity = errors.New("invalid dispatch capacity")
)

// Fulfillment works on Hyperledger smart contract technology.
// It automatically handles verification of dispatch and billing.
type Fulfillment struct {
	contractapi.Contract
}

// TimeWindow is used to describe Dispatch requested time window
type TimeWindow struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (tw TimeWindow) Runtime() int32 {
	return int32((tw.End - tw.Start) / 60)
}

// RequestedDispatch represents a dispatch requested by a consumer
// It is a key structure for verifying agreement fulfillment.
type RequestedDispatch struct {
	ID          string     `json:"id"`
	AgreementID string     `json:"agreement_id"`
	Capacity    float32    `json:"capacity"`
	Status      string     `json:"status"`
	TimeWindow  TimeWindow `json:"time_window"`
	Dispatch    *Dispatch  `json:"dispatch"`
}

// Dispatch represents recorded dispatch of FSP
type Dispatch struct {
	AgreementID string  `json:"agreement_id"`
	RequestID   string  `json:"request_id"`
	Capacity    float32 `json:"capacity"`
	Runtime     int32   `json:"runtime"`
}

// PricingType is used to describe a pricing used in the FSP's bid
type PricingType struct {
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

type Obligation struct {
	// Details
	Capacity   float32       `json:"capacity"`
	MaxRuntime int64         `json:"max_runtime"`
	Pricing    []PricingType `json:"pricing"`

	// Associations
	Asset inventory.Asset `json:"asset"`
}

// Agreement represents an agreed deal between FSP and SO.
// Those BIDS that are in the bidding.Competition and with status TRUE (accepted) are a part of an agreement.
type Agreement struct {
	ID                  string              `json:"id"`
	Competition         bidding.Competition `json:"competition"`
	Obligation          Obligation          `json:"obligation"`
	RequestedDispatches []RequestedDispatch `json:"requested_dispatches"`
}

// RegisterAgreement stores agreement which contains information about the competition
// and its bids. This method should be used after competition auction has finished and
// a set of bids that fulfill competition requirement were accepted.
func (f *Fulfillment) RegisterAgreement(
	ctx contractapi.TransactionContextInterface,
	agreement *Agreement,
) error {
	stub := ctx.GetStub()

	data, err := json.Marshal(agreement)
	if err != nil {
		return err
	}

	return stub.PutState(agreement.ID, data)
}

// RequestDispatch is called by a client when a flexibility consumer request a dispatch
func (f *Fulfillment) RequestDispatch(ctx contractapi.TransactionContextInterface, req RequestedDispatch) error {
	stub := ctx.GetStub()

	data, err := stub.GetState(req.AgreementID)
	if err != nil {
		return err
	}

	agreement := &Agreement{}
	if err := json.Unmarshal(data, agreement); err != nil {
		return nil
	}

	// if dispatch request's capacity is greater than the one agreed upon
	// return an error
	if agreement.Competition.Capacity >= req.Capacity {
		return ErrInvalidDispatchCapacity
	}

	// TODO: validate if request matches competition parameters (e.g service window)

	// Save a dispatch request to the agreement
	agreement.RequestedDispatches = append(agreement.RequestedDispatches, req)

	// Marshal asset object
	asset, err := json.Marshal(agreement.Obligation.Asset)
	if err != nil {
		return fmt.Errorf("failed to emit event: %v", err)
	}

	// Emit event that is used to track metering data
	if err := stub.SetEvent("DispatchRequested", asset); err != nil {
		return err
	}

	data, err = json.Marshal(agreement)
	if err != nil {
		return err
	}

	return stub.PutState(agreement.ID, data)
}

// RegisterDispatch is called after meter data of an asset was requested
func (f *Fulfillment) RegisterDispatch(ctx contractapi.TransactionContextInterface, dispatch Dispatch) error {
	stub := ctx.GetStub()

	data, err := stub.GetState(dispatch.AgreementID)
	if err != nil {
		return err
	}

	agreement := &Agreement{}
	if err := json.Unmarshal(data, agreement); err != nil {
		return nil
	}

	// Find a dispatch request in an agreement
	for _, req := range agreement.RequestedDispatches {
		if req.ID == dispatch.RequestID {
			// TODO: verify runtime and time window

			if req.Capacity < dispatch.Capacity {
				log.Println("requested capacity was not met")
				break
			}

			var invoice string
			for _, price := range agreement.Obligation.Pricing {
				if price.Type == "utilization" {
					invoice = fmt.Sprintf(
						`{"invoice": %.2f}`,
						price.Value*(dispatch.Capacity*float32(dispatch.Runtime/60)),
					)
				}
			}
			if err := stub.SetEvent("DispatchProcessed", []byte(invoice)); err != nil {
				return err
			}

			// Record dispatch
			req.Dispatch = &dispatch
		}
	}

	data, err = json.Marshal(agreement)
	if err != nil {
		return err
	}

	return stub.PutState(agreement.ID, data)
}
