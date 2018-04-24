package main

import (
	"log"
	"time"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
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
	org2User          = "User1"
	channelID         = "orgchannel"
	channelConfigFile = "../../fixtures/fabric/v1.1/channel/orgchannel.tx"
)

//ChainCodeConfig ....
type ChainCodeConfig struct {
	ChainCodeName     string
	ChainCodeVersion  string
	ChainCodePath     string
	ChainCodeGoPath   string
	ChainCodeInitArgs [][]byte
}

func main() {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatal(err)
	}
	defer sdk.Close()

	integration.CleanupUserData(sdk)
	defer integration.CleanupUserData(sdk)

	ccc1 := ChainCodeConfig{
		ChainCodeName:     "examplecc",
		ChainCodeVersion:  "v0",
		ChainCodePath:     "github.com/example_cc",
		ChainCodeGoPath:   "../../chaincode",
		ChainCodeInitArgs: integration.GetInitArgs(),
	}
	log.Printf("chaincode config: %+v\n", ccc1)
	actionWithOrg(sdk, &ccc1)
}

func actionWithOrg(sdk *fabsdk.FabricSDK, ccc *ChainCodeConfig) {
	ordererClientContext := sdk.Context(fabsdk.WithUser(ordererAdminUser), fabsdk.WithOrg(ordererOrgName))
	org1AdminClientContext := sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))
	org2AdminClientContext := sdk.Context(fabsdk.WithUser(org2AdminUser), fabsdk.WithOrg(org2))

	// Channel management client is responsible for managing channels (create/update channel)
	chMgmtClient, err := resmgmt.New(ordererClientContext)
	if err != nil {
		log.Fatal(err)
	}

	// Get signing identity that is used to sign create channel request
	org1AdminUser, err := integration.GetSigningIdentity(sdk, org1AdminUser, org1)
	if err != nil {
		log.Fatalf("failed to get  signing identity of org1AdminUser , err : %v", err)
	}

	org2AdminUser, err := integration.GetSigningIdentity(sdk, org2AdminUser, org2)
	if err != nil {
		log.Fatalf("failed to get  signing identity of org2AdminUser , err : %v", err)
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelConfigFile,
		SigningIdentities: []msp.SigningIdentity{org1AdminUser, org2AdminUser},
	}

	//创建或者更新通道
	txID, err := chMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("SaveChannel failed: %v\n", err)
	}
	log.Printf("create channel txID %s\n", txID.TransactionID)
	//Org1
	org1ResMgmt, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}

	org1Peers, err := getOrgPeers(org1AdminClientContext, org1)
	if err != nil {
		log.Fatalf("get org1 peers failed: %+v\n", err)
	}

	log.Printf("get org1 peers: %+v\n", org1Peers)

	// Org1 peers join channel
	if err = org1ResMgmt.JoinChannel(channelID, resmgmt.WithTargets(org1Peers...), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		log.Fatalf("Org1 peers failed to JoinChannel: %s", err)
	}

	// Org2 resource management client
	org2ResMgmt, err := resmgmt.New(org2AdminClientContext)
	if err != nil {
		log.Fatal(err)
	}

	org2Peers, err := getOrgPeers(org2AdminClientContext, org2)
	if err != nil {
		log.Fatalf("get org2 peers failed: %+v\n", err)
	}

	log.Printf("get org2 peers: %+v\n", org2Peers)
	// Org2 peers join channel // resmgmt.WithTargets(org2Peers...)
	if err = org2ResMgmt.JoinChannel(channelID, resmgmt.WithTargets(org2Peers...), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		log.Fatalf("Org2 peers failed to JoinChannel: %s", err)
	}

	// Create chaincode package
	ccPkg, err := packager.NewCCPackage(ccc.ChainCodePath, ccc.ChainCodeGoPath)
	if err != nil {
		log.Fatal(err)
	}

	installCCReq := resmgmt.InstallCCRequest{Name: ccc.ChainCodeName, Path: ccc.ChainCodePath, Version: ccc.ChainCodeVersion, Package: ccPkg}
	// Install chaincode package to Org1 Peers
	_, err = org1ResMgmt.InstallCC(installCCReq, resmgmt.WithTargets(org1Peers...), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatal(err)
	}

	_, err = org2ResMgmt.InstallCC(installCCReq, resmgmt.WithTargets(org2Peers...), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatal(err)
	}

	//将chaincode策略设置为“两个msps中的任何一个”
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})
	// Org1 resource manager will instantiate chaincode on channel
	instantiateReq := resmgmt.InstantiateCCRequest{
		Name:    ccc.ChainCodeName,
		Path:    ccc.ChainCodePath,
		Version: ccc.ChainCodeVersion,
		Args:    ccc.ChainCodeInitArgs,
		Policy:  ccPolicy,
	}
	instantiateResp, err := org2ResMgmt.InstantiateCC(channelID, instantiateReq, resmgmt.WithTargets(org2Peers[0]), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("Instantiate chaincode err: %v\n", err)
	}
	log.Printf("instantiate response: %+v\n", instantiateResp)

	// check chaincode is instantiated on Org1 peer
	chaincodeQueryResponse, err := org2ResMgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithTargets(org2Peers[0]), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("QueryInstantiatedChaincodes return error: %v", err)
	}
	log.Printf("chaincode query response: %+v\n", chaincodeQueryResponse.Chaincodes)

	// org1AdminChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))
	org2ChannelClientContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org2User), fabsdk.WithOrg(org2))

	// Org2 user connects to 'orgchannel'
	chClientOrg2User, err := channel.New(org2ChannelClientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client for Org2 user: %s", err)
	}
	// Org2 user moves funds on org2 peer
	response, err := chClientOrg2User.Execute(channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetTxArgs()}, channel.WithTargets(org1Peers[1]))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	log.Printf("org2 channel user exec move fonds on org2 peer0 result:%+v\n", response)
	log.Printf("response.Payload: %+v\n", string(response.Payload))
	log.Printf("response.TransactionID: %+v\n", response.TransactionID)
	log.Printf("response.TxValidationCode: %+v\n", response.TxValidationCode)
	log.Println("=======================================================")

	response, err = chClientOrg2User.Query(channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetQueryArgs()}, channel.WithTargets(org2Peers[1]))
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}
	log.Printf("org2 channel user query fonds on org2 peer1 result:%+v\n", response)
	log.Printf("response.Payload: %+v\n", string(response.Payload))
	log.Printf("response.TransactionID: %+v\n", response.TransactionID)
	log.Printf("response.TxValidationCode: %+v\n", response.TxValidationCode)
	log.Println("=======================================================")
	time.Sleep(3 * time.Second)

	org1ChannelClientContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org1User), fabsdk.WithOrg(org1))
	chClientOrg1User, err := channel.New(org1ChannelClientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client for Org1 user: %s", err)
	}
	response, err = chClientOrg1User.Execute(channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetTxArgs()}, channel.WithTargets(org1Peers[0]))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	log.Printf("org1 channel user exec move fonds on org1 peer0 result:%+v\n", response)
	log.Printf("response.Payload: %+v\n", string(response.Payload))
	log.Printf("response.TransactionID: %s\n", response.TransactionID)
	log.Printf("response.TxValidationCode: %+v\n", response.TxValidationCode)
	log.Println("=======================================================")

	response, err = chClientOrg1User.Query(channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetQueryArgs()}, channel.WithTargets(org1Peers[1]))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	log.Printf("org1 user query channel query fonds on org1 peer1 result: %+v\n", response)
	log.Printf("response.Payload: %+v\n", string(response.Payload))
	log.Printf("response.TransactionID: %s\n", response.TransactionID)
	log.Printf("response.TxValidationCode: %+v\n", response.TxValidationCode)
	log.Println("=======================================================")
}

func getOrgPeers(ctxProvider contextAPI.ClientProvider, org string) (peers []fab.Peer, err error) {
	ctx, err := ctxProvider()
	if err != nil {
		return nil, err
	}
	orgPeers, err := ctx.EndpointConfig().PeersConfig(org)
	if err != nil {
		return nil, err
	}
	for _, orgpeer := range orgPeers {
		peer, err := ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: orgpeer})
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	return peers, nil
}
