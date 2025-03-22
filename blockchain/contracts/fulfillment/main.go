package main

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/nikitarudakov/microenergy/pkg/contracts"
	"log"
)

func main() {
	fulfillmentChaincode, err := contractapi.NewChaincode(&contracts.Fulfillment{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := fulfillmentChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
