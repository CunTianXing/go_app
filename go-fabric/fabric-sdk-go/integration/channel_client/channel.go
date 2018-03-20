package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// GenerateRandomID generates random ID
func GenerateRandomID() string {
	return randomString(10)
}

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// FilterTargetsJoinedChannel filters targets to those that have joined the named channel.
func FilterTargetsJoinedChannel(sdk *fabsdk.FabricSDK, orgID string, channelID string, targets []fab.ProposalProcessor) ([]fab.ProposalProcessor, error) {
	joinedTargets := []fab.ProposalProcessor{}
	rc, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgID)).ResourceMgmt()
	if err != nil {
		log.Fatal("failed getting admin user session for org:", err)
	}
	fmt.Println(rc)
	for _, target := range targets {
		// Check if primary peer has joined channel
		alreadyJoined, err := HasPeerJoinedChannel(rc, target, channelID)
		if err != nil {
			log.Fatal(err)
		}
		if alreadyJoined {
			log.Printf("target info: %#v\n", target)
			joinedTargets = append(joinedTargets, target)
		}
	}
	return joinedTargets, nil
}

// HasPeerJoinedChannel checks whether the peer has already joined the channel.
// It returns true if it has, false otherwise, or an error
func HasPeerJoinedChannel(client *resmgmt.Client, peer fab.ProposalProcessor, channel string) (bool, error) {
	foundChannel := false
	response, err := client.QueryChannels(peer)
	if err != nil {
		return false, err
	}
	for _, responseChannel := range response.Channels {
		if responseChannel.ChannelId == channel {
			foundChannel = true
		}
	}
	return foundChannel, nil
}

// CreateChannel attempts to save the named channel.
func CreateChannel(sdk *fabsdk.FabricSDK, req resmgmt.SaveChannelRequest) (bool, error) {
	//Channel management client is responsible for managing channels (create/update)
	resMgmtClient, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg("ordererorg")).ResourceMgmt()
	if err != nil {
		return false, err
	}
	//Create channel (or update if it already exists)
	if err = resMgmtClient.SaveChannel(req); err != nil {
		return false, nil
	}
	time.Sleep(time.Second * 5)
	return true, nil
}

// JoinChannel attempts to save the named channel.
func JoinChannel(sdk *fabsdk.FabricSDK, name string) (bool, error) {
	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := sdk.NewClient(fabsdk.WithUser("Admin")).ResourceMgmt()
	if err != nil {
		log.Println("Failed to create new resource management client")
		return false, err
	}
	if err = resMgmtClient.JoinChannel(name); err != nil {
		return false, nil
	}
	return true, nil
}

// InitializeChannel ...
func InitializeChannel(sdk *fabsdk.FabricSDK, orgID string, req resmgmt.SaveChannelRequest, targets []fab.ProposalProcessor) error {
	joinedTargets, err := FilterTargetsJoinedChannel(sdk, orgID, req.ChannelID, targets)
	if err != nil {
		log.Fatal("checking for joined targets failed:", err)
		return err
	}
	if len(joinedTargets) != len(targets) {
		log.Println("start create channel")
		_, err := CreateChannel(sdk, req)
		if err != nil {
			log.Println("create channel failed")
			return err
		}
		log.Println("end create channel")
		log.Println("start join channel")
		_, err = JoinChannel(sdk, req.ChannelID)
		if err != nil {
			fmt.Println("join channel failed")
			return err
		}
		log.Println("end join channel")
	}
	return nil
}
