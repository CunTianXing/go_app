package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	org1              = "Org1"
	org2              = "Org2"
	org1User          = "User1"
	org1AdminUser     = "Admin"
	org2AdminUser     = "Admin"
	channelID         = "mychannel"
	configFile        = "../../fixtures/config/config_test.yaml"
	channelConfigFile = "../../fixtures/fabric/v1.2/channel/mychannel.tx"
)

func main() {
	baseConfig := integration.BaseConfig{
		ConfigFile:        configFile,
		ChannelID:         channelID,
		OrgID:             org1,
		ChannelConfigFile: channelConfigFile,
	}
	sdk, err := fabsdk.New(config.FromFile(baseConfig.ConfigFile))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}

	if err = baseConfig.Initialize(sdk); err != nil {
		log.Fatalf(err.Error())
	}

	integration.CleanupUserData(sdk)
	defer integration.CleanupUserData(sdk)

	chaincodeID := integration.GenerateRandomID()
	log.Println("chaincodeID:", chaincodeID)
	chaincodePath := "github.com/example_cc"
	chaincodeVersion := "v0"
	pwd, _ := os.Getwd()
	goPath := path.Join(pwd, "../../chaincode")
	fmt.Println(goPath)
	if _, err = integration.InstallAndInstantiateCC(sdk, fabsdk.WithUser(org1AdminUser), baseConfig.OrgID, baseConfig.ChannelID, chaincodeID, chaincodePath, chaincodeVersion, goPath, integration.GetInitArgs()); err != nil {
		log.Fatalf("InstallAndInstantiateCC return error: %v", err)
	}
	log.Println("InstallAndInstantiateCC success!")

	org1AdminChannelContext := sdk.ChannelContext(baseConfig.ChannelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))

	client, err := ledger.New(org1AdminChannelContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}

	ledgerInfo, err := client.QueryInfo()
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}

	configBackend, err := sdk.Config()
	if err != nil {
		log.Fatalf("failed to get config backend, error: %v", err)
	}
	log.Println(configBackend)

	endpointConfig, err := fab.ConfigFromBackend(configBackend)
	if err != nil {
		log.Fatalf("failed to get endpoint config, error: %v", err)
	}
	fmt.Println(endpointConfig)

	expectedPeerConfig, err := endpointConfig.PeerConfig("peer0.org1.example.com")
	if err != nil {
		log.Fatalf("Unable to fetch Peer config for %s", "peer0.org1.example.com")
	}

	if !strings.Contains(ledgerInfo.Endorser, expectedPeerConfig.URL) {
		log.Fatalf("Expecting %s, got %s", expectedPeerConfig.URL, ledgerInfo.Endorser)
	}
	log.Printf("Expecting %s, got %s", expectedPeerConfig.URL, ledgerInfo.Endorser)
	log.Println(ledgerInfo)

	org1ChannelClientContext := sdk.ChannelContext(baseConfig.ChannelID, fabsdk.WithUser(org1User), fabsdk.WithOrg(baseConfig.OrgID))

	chClient, err := channel.New(org1ChannelClientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client: %s", err)
	}
	qreq := channel.Request{ChaincodeID: chaincodeID, Fcn: "invoke", Args: integration.GetQueryArgs()}
	response, err := chClient.Query(qreq)
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}

	value := response.Payload
	log.Println(value)

	treq := channel.Request{ChaincodeID: chaincodeID, Fcn: "invoke", Args: integration.GetTxArgs()}
	log.Printf("req : %+v\n", treq)
	response, err = chClient.Execute(treq)
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}

	response, err = chClient.Query(qreq)
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}

	valueInt, _ := strconv.Atoi(string(value))
	valueAfterInvokeInt, _ := strconv.Atoi(string(response.Payload))
	if valueInt+1 != valueAfterInvokeInt {
		log.Fatalf("Execute failed. Before: %s, after: %s", value, response.Payload)
	}

}
