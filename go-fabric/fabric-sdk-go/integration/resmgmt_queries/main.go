package main

import (
	"log"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/comm"
	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/params"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
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

	chainCodeID := comm.GenerateRandomID()
	if err := comm.InstallAndInstantiateCC(c.SDK, fabsdk.WithUser(c.AdminUser), c.OrgID, channelID, chainCodeID, "github.com/example_cc", "v0", c.GoPath, params.ExampleCCInitArgs()); err != nil {
		log.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
	}
	defer c.SDK.Close()

	client, err := c.SDK.NewClient(fabsdk.WithUser(c.AdminUser)).ResourceMgmt()
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}
	//log.Println(client)
	target := c.Targets[0]
	testInstalledChaincodes(chainCodeID, target, client)
	testInstantiatedChaincodes(c.ChannelID, chainCodeID, target, client)
	testQueryChannels(c.ChannelID, target, client)
	log.Println(target)
}

func testInstalledChaincodes(chainCodeID string, target fab.ProposalProcessor, client *resmgmt.Client) {
	chaincodeQueryResponse, err := client.QueryInstalledChaincodes(target)
	if err != nil {
		log.Fatalf("QueryInstalledChaincodes return error: %v", err)
	}

	found := false
	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		log.Printf("**InstalledCC: %s\n", chaincode)
		if chaincode.Name == chainCodeID {
			found = true
		}
	}

	if !found {
		log.Fatalf("QueryInstalledChaincodes failed to find installed %s chaincode", chainCodeID)
	}
}

func testInstantiatedChaincodes(channelID string, chainCodeID string, target fab.ProposalProcessor, client *resmgmt.Client) {
	chaincodeQueryResponse, err := client.QueryInstantiatedChaincodes(channelID, resmgmt.WithTarget(target.(fab.Peer)))
	if err != nil {
		log.Fatalf("QueryInstantiatedChaincodes return error: %v", err)
	}

	found := false
	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		log.Printf("**InstantiatedCC: %s", chaincode)
		if chaincode.Name == chainCodeID {
			found = true
		}
	}

	if !found {
		log.Fatalf("QueryInstantiatedChaincodes failed to find instantiated %s chaincode", chainCodeID)
	}
}

func testQueryChannels(channelID string, target fab.ProposalProcessor, client *resmgmt.Client) {
	channelQueryResponse, err := client.QueryChannels(target)
	if err != nil {
		log.Fatalf("QueryChannels return error: %v", err)
	}

	found := false
	for _, channel := range channelQueryResponse.Channels {
		log.Printf("**Channel: %s", channel)
		if channel.ChannelId == channelID {
			found = true
		}
	}

	if !found {
		log.Fatalf("QueryChannels failed, peer did not join '%s' channel", channelID)
	}
}
