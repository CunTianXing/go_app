package main

import (
	"log"
	"os"
	"path"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// Initial B values for ExampleCC
const (
	ExampleCCInitA = "100"
	ExampleCCInitB = "200"
)

// ExampleCC query and transaction arguments
var queryArgs = [][]byte{[]byte("query"), []byte("b")}
var txArgs = [][]byte{[]byte("move"), []byte("a"), []byte("b"), []byte("1")}

// ExampleCC init and upgrade args
var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte(ExampleCCInitA), []byte("b"), []byte(ExampleCCInitB)}

// ExampleCCQueryArgs returns example cc query args
func ExampleCCQueryArgs() [][]byte {
	return queryArgs
}

// ExampleCCTxArgs returns example cc move funds args
func ExampleCCTxArgs() [][]byte {
	return txArgs
}

func getChaincodePath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, chainCodePath)
}

//InstallAndInstantiateExampleCC install and instantiate using resource management client
func InstallAndInstantiateExampleCC(sdk *fabsdk.FabricSDK, user fabsdk.IdentityOption, orgName string, chainCodeID string) error {
	return InstallAndInstantiateCC(sdk, user, orgName, chainCodeID, "github.com/example_cc", "v0", getChaincodePath(), initArgs)
}

// InstallAndInstantiateCC install and instantiate using resource management client
func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, user fabsdk.IdentityOption, orgName string, ccName, ccPath, ccVersion, goPath string, ccArgs [][]byte) error {
	ccPkg, err := packager.NewCCPackage(ccPath, goPath)
	if err != nil {
		log.Println("InstallAndInstantiateCC packager.NewCCPackage err:", err)
		return err
	}
	mspID, err := sdk.Config().MspID(orgName)
	if err != nil {
		log.Println("InstallAndInstantiateCC sdk.Config().MspID err:", err)
		return err
	}
	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := sdk.NewClient(user, fabsdk.WithOrg(orgName)).ResourceMgmt()
	if err != nil {
		log.Println("InstallAndInstantiateCC sdk.NewClient err:", err)
		return err
	}
	// install
	_, err = resMgmtClient.InstallCC(resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg})
	if err != nil {
		log.Println("InstallAndInstantiateCC resMgmtClient.InstallCC err:", err)
		return err
	}
	// Instantiate
	ccPolicy := cauthdsl.SignedByMspMember(mspID)
	return resMgmtClient.InstantiateCC(channelID, resmgmt.InstantiateCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Args: ccArgs, Policy: ccPolicy})
}
