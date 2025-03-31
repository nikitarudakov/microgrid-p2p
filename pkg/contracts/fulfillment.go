package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/nikitarudakov/microenergy/internal/onchain"
	"time"
)

type Fulfillment struct {
	contractapi.Contract
}

func (f *Fulfillment) BindContract(ctx contractapi.TransactionContextInterface, contract *onchain.Contract) error {
	stub := ctx.GetStub()

	data, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return stub.PutState(contract.ID, data)
}

func (f *Fulfillment) SetPenalty(ctx contractapi.TransactionContextInterface, contractID string, penalty float64) error {
	stub := ctx.GetStub()

	contract, err := fetchDocByID[onchain.Contract](stub, contractID, "contract")
	if err != nil {
		return err
	}

	contract.Penalty = penalty

	data, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return stub.PutState(contractID, data)
}

// RegisterObligation called after all availability checks were passed
// and a system can start tracking dispatches for an obligation.
func (f *Fulfillment) RegisterObligation(ctx contractapi.TransactionContextInterface, obligation *onchain.Obligation) error {
	stub := ctx.GetStub()

	contract, err := fetchDocByID[onchain.Contract](stub, obligation.ContractID, "contract")
	if err != nil {
		return err
	}

	if err := obligation.Validate(contract); err != nil {
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

	obligation, err := fetchDocByID[onchain.Obligation](stub, obligationID, "obligation")
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

func (f *Fulfillment) RecordDispatch(ctx contractapi.TransactionContextInterface, dispatch *onchain.Dispatch) error {
	stub := ctx.GetStub()

	data, err := json.Marshal(dispatch)
	if err != nil {
		return err
	}

	return stub.PutState(dispatch.ID, data)
}

func (f *Fulfillment) Settlement(ctx contractapi.TransactionContextInterface, dispatchID string) error {
	stub := ctx.GetStub()

	dispatch, err := fetchDocByID[onchain.Dispatch](stub, dispatchID, "dispatch")
	if err != nil {
		return err
	}

	obligation, err := fetchDocByID[onchain.Obligation](stub, dispatch.ObligationID, "obligation")
	if err != nil {
		return err
	}

	contract, err := fetchDocByID[onchain.Contract](stub, obligation.ContractID, "contract")
	if err != nil {
		return err
	}

	settlement := &onchain.Settlement{
		ID:         uuid.New().String(),
		DispatchID: dispatchID,
		SettledAt:  time.Now().Format(time.RFC3339),
		Fiat:       dispatch.CalculatePayableAmount(contract),
	}

	if err := dispatch.Validate(obligation); err != nil {
		settlement.Type = "penalty"
		settlement.Fiat -= contract.Penalty
	} else {
		settlement.Type = "payable"
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
