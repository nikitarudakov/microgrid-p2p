package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/nikitarudakov/microenergy/pkg/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"

	"os"
	"path"
	"time"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/go/src/github.com/nikitarudakov/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	registerArrangement(contract)
	getAllAgreements(contract)
}

func getAllAgreements(contract *client.Contract) {
	evaluateResult, err := contract.EvaluateTransaction("GetAllAgreements")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func registerArrangement(contract *client.Contract) {
	agreement := &contracts.Obligation{
		ID: "Agreement1",
		//Competition: bidding.Competition{
		//	ConsumerName: "ConsumerName",
		//	Capacity:     145.12,
		//	Lifespan:     30,
		//	Voltage: bidding.Voltage{
		//		Min: 2.3,
		//		Max: 4.5,
		//	},
		//	ServiceWindows: []bidding.ServiceWindow{},
		//	Bids:           []bidding.Bid{},
		//},
		//Obligation: contracts.Obligation{
		//	Capacity:   100.00,
		//	MaxRuntime: 10,
		//	Pricing:    []contracts.PricingType{{Type: "utilization", Value: 0.23}},
		//	Asset: inventory.Asset{
		//		MeterID: "M0001",
		//	},
		//},
		//RequestedDispatches: []contracts.RequestedDispatch{},
	}

	data, err := json.Marshal(agreement)
	if err != nil {
		log.Printf("erorr marshalling agreement %+v\n", data)
	}

	print(formatJSON(data))

	_, err = contract.SubmitTransaction("RegisterObligation", formatJSON(data))
	if err != nil {
		var endorseErr *client.EndorseError

		if errors.As(err, &endorseErr) {
			fmt.Printf("Endorse error for transaction %s with gRPC status %v: %s\n", endorseErr.TransactionID, status.Code(endorseErr), endorseErr)
		}
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get home directory: " + err.Error())
	}

	fullTlsCertPath := path.Join(homeDir, cryptoPath, tlsCertPath)

	certificatePEM, err := os.ReadFile(fullTlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get home directory: " + err.Error())
	}

	fullCertPath := path.Join(homeDir, cryptoPath, certPath)

	certificatePEM, err := readFirstFile(fullCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get home directory: " + err.Error())
	}

	fullKeyPath := path.Join(homeDir, cryptoPath, keyPath)

	privateKeyPEM, err := readFirstFile(fullKeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
