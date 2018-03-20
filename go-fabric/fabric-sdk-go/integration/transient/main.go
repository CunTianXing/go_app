package main

import (
	"log"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/comm"
	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/params"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	org1Name          = "Org1"
	configFile        = "../../fixtures/config/config_test0.yaml"
	channelID         = "mychannel"
	channelConfigFile = "../../fixtures/fabric/v1.0.0/channel/mychannel.tx"
	chainCodePath     = "../chaincode"
	adminUser         = "Admin"
)

func main() {
	c := comm.Common{
		ConfigFile:      configFile,
		ChannelID:       channelID,
		ChannelConfig:   channelConfigFile,
		OrgID:           org1Name,
		AdminUser:       adminUser,
		GoPath:          chainCodePath,
		ConnectEventHub: true,
	}

	if err := c.Initialize(); err != nil {
		log.Fatalf(err.Error())
	}
	defer c.SDK.Close()
	chainCodeID := comm.GenerateRandomID()
	if err := comm.InstallAndInstantiateCC(c.SDK, fabsdk.WithUser(c.AdminUser), c.OrgID, c.ChannelID, chainCodeID, "github.com/example_cc", "v0", c.GoPath, params.ExampleCCInitArgs()); err != nil {
		log.Fatalf("InstallAndInstantiateCC return error: %v", err)
	}

	fcn := "invoke"
	transientData := "Transient data test..."
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte(transientData)
	//(transactor fab.ProposalSender, chainCodeID string, fcn string, args [][]byte, targets []fab.ProposalProcessor, transientData map[string][]byte)
	transactionProposalResponse, _, err := comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, params.ExampleCCTxArgs(), c.Targets[:1], transientDataMap)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal return error: %v", err)
	}
	log.Printf("CreateAndSendTransactionProposal result: %#v\n", transactionProposalResponse)
	strResponse := string(transactionProposalResponse[0].ProposalResponse.GetResponse().Payload)
	//validate transient data exists in proposal
	if len(strResponse) == 0 {
		log.Fatalf("Transient data does not exist: expected %s", transientData)
	}
	log.Printf("strResponse:%s\n", strResponse)
	//verify transient data content
	if strResponse != transientData {
		log.Fatalf("Expected '%s' in transient data field. Received '%s' ", transientData, strResponse)
	}
	//transient data null
	transientDataMap["result"] = []byte{}
	transactionProposalResponse, _, err = comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, params.ExampleCCTxArgs(), c.Targets[:1], transientDataMap)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal with empty transient data return an error: %v", err)
	}
	//validate that transient data does not exist in proposal
	strResponse = string(transactionProposalResponse[0].ProposalResponse.GetResponse().Payload)
	if len(strResponse) != 0 {
		log.Fatalf("Transient data validation has failed. An empty transient data was expected but %s was returned", strResponse)
	}
	log.Printf("strResponse:%s\n", strResponse)
}
