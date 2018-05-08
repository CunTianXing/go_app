package integration

import (
	"log"
	"math/rand"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	pmsp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"github.com/pkg/errors"
)

const (
	adminUser      = "Admin"
	ordererOrgName = "ordererorg"
)

// GenerateRandomID generates random ID
func GenerateRandomID() string {
	return randomString(10)
}

// Utility to create random string of strlen length
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// CleanupPath removes the contents of a state store.
func CleanupPath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		log.Fatal(err)
	}
}

//CleanupUserData removes user data.
func CleanupUserData(sdk *fabsdk.FabricSDK) {
	configBackend, err := sdk.Config()
	if err != nil {
		log.Fatal(err)
	}

	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	identityConfig, err := msp.ConfigFromBackend(configBackend)
	if err != nil {
		log.Fatal(err)
	}

	keyStorePath := cryptoSuiteConfig.KeyStorePath()
	credentialStorePath := identityConfig.CredentialStorePath()
	CleanupPath(keyStorePath)
	CleanupPath(credentialStorePath)
}

//GetSigningIdentity ....
func GetSigningIdentity(sdk *fabsdk.FabricSDK, user, org string) (pmsp.SigningIdentity, error) {
	identityContext := sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(org))
	return identityContext()
}

// InitializeChannel ...
func InitializeChannel(sdk *fabsdk.FabricSDK, orgID string, req resmgmt.SaveChannelRequest, targets []string) error {
	joinedTargets, err := FilterTargetsJoinedChannel(sdk, orgID, req.ChannelID, targets)
	if err != nil {
		log.Printf("checking for joined targets failed : %v\n", err)
		return err
	}
	if len(joinedTargets) != len(targets) {
		_, err := CreateChannel(sdk, req)
		if err != nil {
			log.Printf("create channel failed")
			return err
		}
		_, err = JoinChannel(sdk, req.ChannelID, orgID)
		if err != nil {
			log.Printf("join channel failed")
			return err
		}
	}
	return nil
}

//CreateChannel ...
func CreateChannel(sdk *fabsdk.FabricSDK, req resmgmt.SaveChannelRequest) (bool, error) {
	clientContext := sdk.Context(fabsdk.WithUser(adminUser), fabsdk.WithOrg(ordererOrgName))

	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Printf("Failed to create new channel management client: %v\n", err)
		return false, err
	}

	if _, err = resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		return false, err
	}
	return true, nil
}

//JoinChannel ....
func JoinChannel(sdk *fabsdk.FabricSDK, channelID, orgID string) (bool, error) {
	clientContext := sdk.Context(fabsdk.WithUser(adminUser), fabsdk.WithOrg(orgID))

	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Printf("Failed to create new resource management client:%v\n", err)
		return false, err
	}

	if err = resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		return false, err
	}
	return true, nil
}

//FilterTargetsJoinedChannel ....
func FilterTargetsJoinedChannel(sdk *fabsdk.FabricSDK, orgID string, channelID string, targets []string) ([]string, error) {
	var joinedTargets []string

	//prepare context
	clientContext := sdk.Context(fabsdk.WithUser(adminUser), fabsdk.WithOrg(orgID))

	rc, err := resmgmt.New(clientContext)
	if err != nil {
		return nil, errors.WithMessage(err, "failed getting admin user session for org")
	}

	for _, target := range targets {
		// Check if primary peer has joined channel
		alreadyJoined, err := HasPeerJoinedChannel(rc, target, channelID)
		if err != nil {
			return nil, errors.WithMessage(err, "failed while checking if primary peer has already joined channel")
		}
		if alreadyJoined {
			joinedTargets = append(joinedTargets, target)
		}
	}
	return joinedTargets, nil
}

//HasPeerJoinedChannel ....
func HasPeerJoinedChannel(client *resmgmt.Client, target string, channel string) (bool, error) {
	foundChannel := false
	response, err := client.QueryChannels(resmgmt.WithTargetURLs(target))
	if err != nil {
		log.Printf("failed to query channel for peer: %v\n", err)
		return false, err
	}
	for _, responseChannel := range response.Channels {
		if responseChannel.ChannelId == channel {
			foundChannel = true
		}
	}
	return foundChannel, nil
}
