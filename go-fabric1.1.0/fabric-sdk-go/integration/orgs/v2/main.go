package main

import (
	"log"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

const (
	configFile       = "../../../fixtures/config/config_test.yaml"
	ordererAdminUser = "Admin"
	ordererOrgName   = "ordererorg"
	org1AdminUser    = "Admin"
	org2AdminUser    = "Admin"
	pollRetries      = 5
	org1             = "Org1"
	org2             = "Org2"
	org1User         = "User1"
	org2User         = "User1"
)

//ChainCodeConfig ....
type ChainCodeConfig struct {
	ChainCodeName     string
	ChainCodeVersion  string
	ChainCodePath     string
	ChainCodeGoPath   string
	ChainCodeInitArgs [][]byte
}

//ChannelConfig ...
type ChannelConfig struct {
	ChannelID         string
	ChannelConfigPath string
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
		ChainCodeGoPath:   "../../../chaincode",
		ChainCodeInitArgs: integration.GetInitArgs(),
	}
	ccc2 := ChainCodeConfig{
		ChainCodeName:     "examplecc",
		ChainCodeVersion:  "v1",
		ChainCodePath:     "github.com/example_cc",
		ChainCodeGoPath:   "../../../chaincode",
		ChainCodeInitArgs: integration.GetInitArgs(),
	}
	log.Printf("chaincode config: %+v\n", ccc1)
	chc1 := ChannelConfig{
		ChannelID:         "orgchannel",
		ChannelConfigPath: "../../../fixtures/fabric/v1.1/channel/orgchannel.tx",
	}
	log.Printf("channel config: %+v\n", chc1)
	actionWithOrg(sdk, &ccc1, &ccc2, &chc1)
}

func actionWithOrg(sdk *fabsdk.FabricSDK, ccc *ChainCodeConfig, ccc2 *ChainCodeConfig, chc *ChannelConfig) {
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

	//创建或者更新通道
	signingIdentities := []msp.SigningIdentity{org1AdminUser, org2AdminUser}
	createChannel(chMgmtClient, signingIdentities, chc)
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
	if err = org1ResMgmt.JoinChannel(chc.ChannelID, resmgmt.WithTargets(org1Peers...), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
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
	if err = org2ResMgmt.JoinChannel(chc.ChannelID, resmgmt.WithTargets(org2Peers...), resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
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
	instantiateResp, err := org2ResMgmt.InstantiateCC(chc.ChannelID, instantiateReq, resmgmt.WithTargets(org2Peers[0]), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("Instantiate chaincode err: %v\n", err)
	}
	log.Printf("instantiate response: %+v\n", instantiateResp)

	// check chaincode is instantiated on Org1 peer
	chaincodeQueryResponse, err := org2ResMgmt.QueryInstantiatedChaincodes(chc.ChannelID, resmgmt.WithTargets(org2Peers[0]), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("QueryInstantiatedChaincodes return error: %v", err)
	}
	log.Printf("chaincode query response: %+v\n", chaincodeQueryResponse.Chaincodes)

	// org1AdminChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))
	org2ChannelClientContext := sdk.ChannelContext(chc.ChannelID, fabsdk.WithUser(org2User), fabsdk.WithOrg(org2))

	// Org2 user connects to 'orgchannel'
	chClientOrg2User, err := channel.New(org2ChannelClientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client for Org2 user: %s", err)
	}
	// Org2 user moves funds on org1 peer1

	log.Println("org2 channel user exec move fonds on org1 peer1")
	moveFunds(chClientOrg2User, ccc, org1Peers[1])

	//query
	log.Println("org2 channel user query fonds on org2 peer1")
	query(chClientOrg2User, ccc, org2Peers[1])

	org1ChannelClientContext := sdk.ChannelContext(chc.ChannelID, fabsdk.WithUser(org1User), fabsdk.WithOrg(org1))
	chClientOrg1User, err := channel.New(org1ChannelClientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client for Org1 user: %s", err)
	}

	log.Println("org1 channel user exec move fonds on org1 peer0")
	moveFunds(chClientOrg1User, ccc, org1Peers[0])

	//query
	log.Println("org1 user query channel query fonds on org1 peer1 ")
	query(chClientOrg1User, ccc, org1Peers[1])

	//upgrade
	resmgmt := []resmgmt.Client{*org1ResMgmt, *org2ResMgmt}
	upgradeCC(resmgmt, ccPkg, ccc2, chc)
	checkChaincodePolicy(chClientOrg2User, ccc, org2Peers, org1Peers)
}

func createChannel(chMgmtClient *resmgmt.Client, orgAdminUsers []msp.SigningIdentity, chc *ChannelConfig) {
	req := resmgmt.SaveChannelRequest{ChannelID: chc.ChannelID, ChannelConfigPath: chc.ChannelConfigPath, SigningIdentities: orgAdminUsers}
	saveChannelResponse, err := chMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("savechannel failed: %s\n", err.Error())
	}
	log.Printf("create channel txID %s\n", saveChannelResponse.TransactionID)
}

func query(chClientOrgUser *channel.Client, ccc *ChainCodeConfig, target fab.Peer) []byte {
	request := channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetQueryArgs()}
	response, err := chClientOrgUser.Query(request, channel.WithTargets(target), channel.WithRetry(retry.DefaultChClientOpts))
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}
	log.Printf("response.Payload: %+v\n", string(response.Payload))
	log.Printf("response.TransactionID: %s\n", response.TransactionID)
	log.Printf("response.TxValidationCode: %+v\n", response.TxValidationCode)
	return response.Payload
}

func moveFunds(chClientOrgUser *channel.Client, ccc *ChainCodeConfig, target fab.Peer) fab.TransactionID {
	request := channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetTxArgs()}
	response, err := chClientOrgUser.Execute(request, channel.WithTargets(target), channel.WithRetry(retry.DefaultChClientOpts))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}

	if response.ChaincodeStatus == 0 {
		log.Fatalf("Expected ChaincodeStatus:%d\n ", response.ChaincodeStatus)
	}
	if response.Responses[0].ChaincodeStatus != response.ChaincodeStatus {
		log.Fatalf("Expected the chaincode status returned by successful Peer Endorsement to be same as Chaincode status for client response")
	}

	log.Printf("response.Payload: %+v\n", string(response.Payload))
	log.Printf("response.TransactionID: %s\n", response.TransactionID)
	log.Printf("response.TxValidationCode: %+v\n", response.TxValidationCode)
	log.Println("=======================================================")
	return response.TransactionID
}

func upgradeCC(orgRMgmtClients []resmgmt.Client, ccPkg *api.CCPackage, ccc *ChainCodeConfig, chc *ChannelConfig) {
	request := resmgmt.InstallCCRequest{Name: ccc.ChainCodeName, Path: ccc.ChainCodePath, Version: ccc.ChainCodeVersion, Package: ccPkg}
	for _, client := range orgRMgmtClients {
		_, err := client.InstallCC(request, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			log.Fatalf("Failed install chaincode on %s : %v\n: ", client, err)
		}
	}
	// New chaincode policy
	org1Andorg2Policy, err := cauthdsl.FromString("AND ('Org1MSP.member','Org2MSP.member')")
	if err != nil {
		log.Fatalf("Failed new chaincode policy: %v\n", err)
	}
	upgradeReq := resmgmt.UpgradeCCRequest{Name: ccc.ChainCodeName, Path: ccc.ChainCodePath, Version: ccc.ChainCodeVersion, Args: integration.GetUpgradeArgs(), Policy: org1Andorg2Policy}
	upgradeResp, err := orgRMgmtClients[0].UpgradeCC(chc.ChannelID, upgradeReq)
	if err != nil {
		log.Fatalf("upgrade chaincode failed: %s\n", err)
	}
	log.Printf("upgrade chaincode response txID: %s\n", upgradeResp.TransactionID)
}

func checkChaincodePolicy(chOrg2Client *channel.Client, ccc *ChainCodeConfig, org1Peers []fab.Peer, org2Peers []fab.Peer) {
	request := channel.Request{ChaincodeID: ccc.ChainCodeName, Fcn: "invoke", Args: integration.GetTxArgs()}
	response, err := chOrg2Client.Execute(request, channel.WithTargets(org1Peers[0], org2Peers[0]), channel.WithRetry(retry.DefaultChClientOpts))
	if err != nil {
		log.Printf("Should have failed to move funds due to cc policy: [%s %s]\n", org1Peers[0], org2Peers[0])
	} else {
		log.Printf("org channel user exec move fonds on org peer [%s %s] response txID: %s\n", org1Peers[0], org2Peers[0], response.TransactionID)
	}

	response, err = chOrg2Client.Execute(request, channel.WithTargets(org1Peers...), channel.WithRetry(retry.DefaultChClientOpts))
	if err != nil {
		log.Fatalf("Should have failed to move funds due to cc policy: %v\n", org1Peers)
	}
	log.Printf("org channel user exec move fonds on org peer %v response txID: %s\n", org1Peers, response.TransactionID)
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
