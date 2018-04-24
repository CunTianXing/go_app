package main

import (
	"log"
	"os"
	"time"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"google.golang.org/grpc/grpclog"
)

var (
	org1              = "Org1"
	org2              = "Org2"
	ordererAdminUser  = "Admin"
	ordererOrgName    = "ordererorg"
	org1AdminUser     = "Admin"
	org2AdminUser     = "Admin"
	logger            = logging.NewLogger("test-logger")
	configFile        = "../../fixtures/config/config_expired_orderers_cert_test.yaml"
	channelConfigFile = "../../fixtures/fabric/v1.2/channel/orgchannel.tx"
)

func main() {
	os.Setenv("GRPC_TRACE", "all")
	os.Setenv("GRPC_VERBOSITY", "DEBUG")
	os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "INFO")
	grpclog.SetLogger(logger)
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatal(err)
	}
	defer sdk.Close()
	time.Sleep(100 * time.Microsecond)

	integration.CleanupUserData(sdk)
	defer integration.CleanupUserData(sdk)

	ordererClientContext := sdk.Context(fabsdk.WithUser(ordererAdminUser), fabsdk.WithOrg(ordererOrgName))

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
	_, err = chMgmtClient.SaveChannel(req)
	if err != nil {
		log.Fatalf("Expected error: calling orderer 'orderer.example.com:7050' failed: Orderer Client Status Code: (2) CONNECTION_FAILED....%s\n", err)
	}
	time.Sleep(100 * time.Millisecond)
}
