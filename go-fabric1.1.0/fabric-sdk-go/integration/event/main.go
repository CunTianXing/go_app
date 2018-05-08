package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	configFile        = "../../fixtures/config/config_test.yaml"
	ordererAdminUser  = "Admin"
	ordererOrgName    = "ordererorg"
	org1AdminUser     = "Admin"
	org2AdminUser     = "Admin"
	pollRetries       = 5
	org1              = "Org1"
	org2              = "Org2"
	org1User          = "User1"
	org2User          = "User2"
	channelID         = "mychannel"
	channelConfigFile = "../../fixtures/fabric/v1.1/channel/mychannel.tx"
)

func main() {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()

	integration.CleanupUserData(sdk)
	defer integration.CleanupUserData(sdk)

	baseConfig := integration.BaseConfig{
		ConfigFile:        configFile,
		ChannelID:         channelID,
		OrgID:             org1,
		ChannelConfigFile: channelConfigFile,
	}

	if err = baseConfig.Initialize(sdk); err != nil {
		log.Fatalf(err.Error())
	}

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

	org1ChannelClientContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org1User), fabsdk.WithOrg(org1))
	chClient, err := channel.New(org1ChannelClientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client: %s", err)
	}

	eventClient, err := event.New(org1ChannelClientContext, event.WithBlockEvents())
	if err != nil {
		log.Fatalf("Failed to create new events client: %s", err)
	}

	actionEvent(chaincodeID, chClient, eventClient, true)

	RegisterFilteredBlockEvent(chaincodeID, chClient, eventClient)

	_, _, err = eventClient.RegisterBlockEvent()
	if err != nil {
		log.Fatalf("Default events client should have failed to register for block events")
	}
}

func actionEvent(ccID string, chClient *channel.Client, eventClient *event.Client, expectPayload bool) {
	eventID := integration.GenerateRandomID()
	payload := "Test Payload"

	reg, notifier, err := eventClient.RegisterChaincodeEvent(ccID, eventID)
	if err != nil {
		log.Fatalf("Failed to register cc event: %s", err)
	}
	defer eventClient.Unregister(reg)

	response, err := chClient.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: append(integration.GetTxArgs(), []byte(eventID))})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}

	select {
	case ccEvent := <-notifier:
		log.Printf("Received cc event.EventName: %s\n", ccEvent.EventName)
		log.Printf("Received cc event.BlockNumber: %d\n", ccEvent.BlockNumber)
		log.Printf("Received cc event.ChaincodeID: %s\n", ccEvent.ChaincodeID)
		log.Printf("Received cc event.TxID: %s\n", ccEvent.TxID)
		log.Printf("Received cc event.SourceURL: %s\n", ccEvent.SourceURL)
		log.Printf("Received cc event.Payload: %s\n", string(ccEvent.Payload))
		if expectPayload && string(ccEvent.Payload[:]) != payload {
			log.Fatalf("Did not receive 'Test Payload'")
		}

		if !expectPayload && string(ccEvent.Payload[:]) != "" {
			log.Fatalf("Expected empty payload, got %s", ccEvent.Payload[:])
		}
		if ccEvent.TxID != string(response.TransactionID) {
			log.Fatalf("CCEvent(%s) and Execute(%s) transaction IDs don't match", ccEvent.TxID, string(response.TransactionID))
		}

	case <-time.After(time.Second * 20):
		log.Fatalf("Did NOT receive CC for eventId(%s)\n", eventID)
	}
}

func RegisterBlockEvent(ccID string, chClient *channel.Client, eventClient *event.Client) {
	breg, beventch, err := eventClient.RegisterBlockEvent()
	if err != nil {
		log.Fatalf("Error registering for block events: %s", err)
	}
	defer eventClient.Unregister(breg)

	response, err := chClient.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: integration.GetTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	log.Printf("response cc TxID: %s\n", response.TransactionID)
	select {
	case e, ok := <-beventch:
		if !ok {
			log.Fatalf("unexpected closed channel while waiting for block event")
		}
		log.Printf("Received block event: %#v", e)
		if e.Block == nil {
			log.Fatalf("Expecting block in block event but got nil")
		}
	case <-time.After(time.Second * 20):
		log.Fatalf("Did NOT receive block event for txID(%s)\n", response.TransactionID)
	}
}

func RegisterFilteredBlockEvent(ccID string, chClient *channel.Client, eventClient *event.Client) {
	fbreg, fbeventch, err := eventClient.RegisterFilteredBlockEvent()
	if err != nil {
		log.Fatalf("Error registering for block events: %s", err)
	}
	defer eventClient.Unregister(fbreg)

	response, err := chClient.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: integration.GetTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	log.Printf("response cc TxID: %s\n", response.TransactionID)
	select {
	case event, ok := <-fbeventch:
		if !ok {
			log.Fatalf("unexpected closed channel while waiting for filtered block event")
		}
		if event.FilteredBlock == nil {
			log.Fatalf("Expecting filtered block in filtered block event but got nil")
		}
		log.Printf("Received filtered block event: %#v", event)
	case <-time.After(time.Second * 20):
		log.Fatalf("Did NOT receive filtered block event for txID(%s)\n", response.TransactionID)
	}
}
