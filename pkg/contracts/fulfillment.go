package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"math"
	"time"
)

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
	MaxRuntimeSeconds int32   `json:"max_runtime_seconds"`
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

func (o Obligation) validate(c *Contract) error {
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
	ID            string     `json:"id"`
	DocType       string     `json:"doc_type"`
	ObligationID  string     `json:"obligation_id"`
	TimeWindow    TimeWindow `json:"time_window"`
	Direction     string     `json:"direction"`
	Capacity      float64    `json:"capacity"`
	TimeAvailable int64      `json:"time_available"`
}

func (d *Dispatch) calculatePayableAmount(c *Contract, o *Obligation) float64 {
	if c.ServiceFee != 0 {
		return c.ServiceFee
	}

	var totalAmount float64
	hours, _ := d.TimeWindow.Hours()

	// Add availability amount
	totalAmount += c.Capacity * c.AvailabilityPrice * hours

	// Add dispatch amount
	totalAmount += math.Abs(o.TimeWindow.Baseline-d.Capacity) * c.UtilizationPrice

	return totalAmount
}

func (d *Dispatch) validate(o *Obligation) error {
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

type Fulfillment struct {
	contractapi.Contract
}

func (f *Fulfillment) BindContract(ctx contractapi.TransactionContextInterface, contract *Contract) error {
	stub := ctx.GetStub()

	data, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return stub.PutState(contract.ID, data)
}

// RegisterObligation called after all availability checks were passed
// and a system can start tracking dispatches for an obligation.
func (f *Fulfillment) RegisterObligation(ctx contractapi.TransactionContextInterface, obligation *Obligation) error {
	stub := ctx.GetStub()

	contract, err := fetchDocByID[Contract](stub, obligation.ContractID, "contract")
	if err != nil {
		return err
	}

	if err := obligation.validate(contract); err != nil {
		return err
	}

	data, err := json.Marshal(obligation)
	if err != nil {
		return err
	}

	return stub.PutState(obligation.ID, data)
}

// RecordObligationStoppage should result in immediate settlement
func (f *Fulfillment) RecordObligationStoppage(ctx contractapi.TransactionContextInterface, obligationID string) error {
	stub := ctx.GetStub()

	obligation, err := fetchDocByID[Obligation](stub, obligationID, "obligation")
	if err != nil {
		return err
	}

	obligation.Status = "stop_requested"

	data, err := json.Marshal(obligation)
	if err != nil {
		return err
	}

	return stub.PutState(obligation.ID, data)
}

func (f *Fulfillment) RecordDispatch(ctx contractapi.TransactionContextInterface, dispatch *Dispatch) error {
	stub := ctx.GetStub()

	data, err := json.Marshal(dispatch)
	if err != nil {
		return err
	}

	return stub.PutState(dispatch.ID, data)
}

func (f *Fulfillment) Settlement(ctx contractapi.TransactionContextInterface, dispatchID string) error {
	stub := ctx.GetStub()

	dispatch, err := fetchDocByID[Dispatch](stub, dispatchID, "dispatch")
	if err != nil {
		return err
	}

	obligation, err := fetchDocByID[Obligation](stub, dispatch.ObligationID, "obligation")
	if err != nil {
		return err
	}

	contract, err := fetchDocByID[Contract](stub, obligation.ContractID, "contract")
	if err != nil {
		return err
	}

	settlement := &Settlement{
		ID:         uuid.New().String(),
		DispatchID: dispatchID,
		SettledAt:  time.Now().Format(time.RFC3339),
	}

	if err := dispatch.validate(obligation); err != nil {
		settlement.Type = "penalty"
		settlement.Fiat = contract.Penalty
	} else {
		settlement.Type = "payable"
		settlement.Fiat = dispatch.calculatePayableAmount(contract, obligation)
	}

	data, err := json.Marshal(settlement)
	if err != nil {
		return err
	}

	return stub.PutState(settlement.ID, data)
}

func fetchDocByID[T interface{}](stub shim.ChaincodeStubInterface, id, docType string) (*T, error) {
	query := fmt.Sprintf(`{
		"selector": {
			"doc_type": "%s",
			"id": "%s"
		}
	}`, docType, id)

	results, err := stub.GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	if results.HasNext() {
		res, err := results.Next()
		if err != nil {
			return nil, err
		}

		var o T
		if err := json.Unmarshal(res.Value, &o); err != nil {
			return nil, err
		}

		return &o, nil
	}

	return nil, errors.New("not found")
}
