package main

import (
	"fmt"
	"log"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	org1              = "Org1"
	org2              = "Org2"
	ordererAdminUser  = "Admin"
	ordererOrgName    = "ordererorg"
	org1AdminUser     = "Admin"
	org2AdminUser     = "Admin"
	configFile        = "../../fixtures/config/config_expired_peers_cert_test.yaml"
	channelConfigFile = "../../fixtures/fabric/v1.2/channel/orgchannel.tx"
)

func main() {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatal(err)
	}
	defer sdk.Close()

	integration.CleanupUserData(sdk)
	defer integration.CleanupUserData(sdk)

	ordererClientContext := sdk.Context(fabsdk.WithUser(ordererAdminUser), fabsdk.WithOrg(ordererOrgName))
	org1AdminClientContext := sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))

	chMgmtClient, err := resmgmt.New(ordererClientContext)
	if err != nil {
		log.Fatal(err)
	}

	org1AdminUser, err := integration.GetSigningIdentity(sdk, org1AdminUser, org1)
	if err != nil {
		log.Fatal(err)
	}

	org2AdminUser, err := integration.GetSigningIdentity(sdk, org2AdminUser, org2)
	if err != nil {
		log.Fatal(err)
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         "orgchannel",
		ChannelConfigPath: channelConfigFile,
		SigningIdentities: []msp.SigningIdentity{org1AdminUser, org2AdminUser},
	}

	txID, err := chMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(txID)

	org1ResMgmt, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		log.Fatal(err)
	}

	err = org1ResMgmt.JoinChannel("orgchannel", resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println("dddd")
		log.Fatal(err)
	}

}
