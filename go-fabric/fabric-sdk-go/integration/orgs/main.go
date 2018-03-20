package main

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	selection "github.com/hyperledger/fabric-sdk-go/pkg/client/common/selection/dynamicselection"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defsvc"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

const (
	pollRetries = 5
	org1        = "Org1"
	org2        = "Org2"
)

// Peers
var orgTestPeer0 fab.Peer
var orgTestPeer1 fab.Peer

// config
var (
	channelID         = "orgchannel"
	channelConfigPath = "../../fixtures/fabric/v1.0.0/channel/orgchannel.tx"
	configFile        = "../../fixtures/config/config_test.yaml"
	chainCodePath     = "../chaincode"
)

//创建一个有两个组织的通道，安装chaincode在它们中的每一个上，最后在org2对等节点上调用事务并查询来自org1对等的结果
func main() {
	// Create SDK setup
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()
	expectedValue := testWithOrg1(sdk)
	expectedValue = testWithOrg2(expectedValue)
	verifyWithOrg1(sdk, expectedValue)
}

func testWithOrg1(sdk *fabsdk.FabricSDK) int {
	// 通道管理客户端负责管理管理（创建/更新通道）
	chMgmtClient, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg("ordererorg")).ResourceMgmt()
	if err != nil {
		log.Fatalf("Failed to new channel management client: %s", err)
	}

	// Create channel (or update if it already exists)
	org1AdminUser := loadOrgUser(sdk, org1, "Admin")
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfig: channelConfigPath, SigningIdentity: org1AdminUser}
	if err = chMgmtClient.SaveChannel(req); err != nil {
		log.Fatalf("failed save channel: %s", err)
	}
	//Allow orderer to process channel creation
	time.Sleep(time.Second * 5)

	// Org1资源管理客户端（Org1是默认组织）
	org1ResMgmt, err := sdk.NewClient(fabsdk.WithUser("Admin")).ResourceMgmt()
	if err != nil {
		log.Fatalf("Failed to create new org1 resource management client: %s", err)
	}

	// Org1 peers join channel
	if err = org1ResMgmt.JoinChannel(channelID); err != nil {
		log.Fatalf("Org1 peers failed to JoinChannel: %s", err)
	}
	log.Printf("Org1 peers join channel: %s\n", channelID)

	// Org2 resource management client
	org2ResMgmt, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg(org2)).ResourceMgmt()
	if err != nil {
		log.Fatalf("Failed to create new org2 resource management client: %s", err)
	}

	// Org2 peers join channel
	if err = org2ResMgmt.JoinChannel(channelID); err != nil {
		log.Fatalf("Org2 peers failed to JoinChannel: %s", err)
	}
	log.Printf("Org2 peers join channel: %s\n", channelID)

	// Create chaincode package for example cc
	ccPkg, err := packager.NewCCPackage("github.com/example_cc", chainCodePath)
	if err != nil {
		log.Fatalf("Failed to create chaincode package: %s", err)
	}

	installCCReq := resmgmt.InstallCCRequest{Name: "exampleCC", Path: "github.com/example_cc", Version: "0", Package: ccPkg}
	// Install example cc to Org1 peers
	response1, err := org1ResMgmt.InstallCC(installCCReq)
	if err != nil {
		log.Fatalf("Failed to install example cc to Org1 peers: %s", err)
	}
	log.Printf("Install example cc to Org1 peers response: %#v\n", response1)

	// Install example cc to Org2 peers
	response2, err := org2ResMgmt.InstallCC(installCCReq)
	if err != nil {
		log.Fatalf("Failed to install example cc to Org2 peers: %s", err)
	}
	log.Printf("Install example cc to Org2 peers response: %#v\n", response2)
	log.Println("================Install chaincode ok====================")
	log.Println("Set up chaincode policy to 'any of two msps'")
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})
	// Org1 资源管理者将实例化 'example_cc' 在 'orgchannel'
	if err = org1ResMgmt.InstantiateCC(channelID, resmgmt.InstantiateCCRequest{Name: "exampleCC", Path: "github.com/example_cc", Version: "0", Args: ExampleCCInitArgs(), Policy: ccPolicy}); err != nil {
		log.Fatalf("Failed to org1 instantiate example_cc: %s", err)
	}
	log.Println("org1 instantiate success!")
	//同一通道上chaincode 只能实例化一次

	//Load specific targets for move funds test
	loadOrgPeers(sdk)

	//Verify that example CC is instantiated on Org1 peer
	chaincodeQueryResponse, err := org1ResMgmt.QueryInstantiatedChaincodes(channelID)
	if err != nil {
		log.Fatalf("QueryInstantiatedChaincodes return error: %v", err)
	}
	log.Printf("all chaincodes info:%#v\n", chaincodeQueryResponse)
	found := false
	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		log.Printf("chaincode info: %#v\n", chaincode)
		if chaincode.Name == "exampleCC" {
			found = true
		}
	}
	if !found {
		log.Fatalf("QueryInstantiatedChaincodes failed to find instantiated exampleCC chaincode")
	}

	// Org1 user connects to 'orgchannel'
	chClientOrg1User, err := sdk.NewClient(fabsdk.WithUser("User1"), fabsdk.WithOrg(org1)).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to org1 user connect to orgchannel: %s", err)
	}
	log.Printf("Org1 user connects to '%s' success!", channelID)
	// Org2 user connects to 'orgchannel'

	chClientOrg2User, err := sdk.NewClient(fabsdk.WithUser("User1"), fabsdk.WithOrg(org2)).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to org2 user connect to orgchannel: %s", err)
	}
	log.Printf("Org2 user connects to '%s' success!", channelID)
	log.Println(chClientOrg2User)
	// Org1用户查询两个对等节点的初始值
	response, err := chClientOrg1User.Query(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: ExampleCCQueryArgs()})
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}
	log.Printf("org1 user query initial value on both peers: %#v\n", response)
	initial, _ := strconv.Atoi(string(response.Payload))
	log.Println("Payload:", initial)

	// Ledger客户端将验证区块链信息
	ledgerClient, err := sdk.NewClient(fabsdk.WithUser("Admin")).Ledger(channelID)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}
	log.Printf("create new ledger client success!")

	channelCfg, err := ledgerClient.QueryConfig(ledger.WithTargets(orgTestPeer0.(fab.Peer), orgTestPeer1.(fab.Peer)), ledger.WithMinTargets(2))
	if err != nil {
		log.Fatalf("QueryConfig return error: %v", err)
	}
	log.Printf("QueryConfig return cfgInfo: %#v\n", channelCfg)

	if len(channelCfg.Orderers()) == 0 {
		log.Fatalf("Failed to retrieve channel orderers")
	}
	expectedOrderer := "orderer.example.com"
	if !strings.Contains(channelCfg.Orderers()[0], expectedOrderer) {
		log.Fatalf("Expecting %s, got %s", expectedOrderer, channelCfg.Orderers()[0])
	}

	ledgerInfoBefore, err := ledgerClient.QueryInfo(ledger.WithTargets(orgTestPeer0.(fab.Peer), orgTestPeer1.(fab.Peer)), ledger.WithMinTargets(2), ledger.WithMaxTargets(3))
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}
	log.Printf("QueryInfo return ledgerInfo: %#v\n", ledgerInfoBefore)

	// Test Query Block by Hash - retrieve current block by hash
	block, err := ledgerClient.QueryBlockByHash(ledgerInfoBefore.BCI.CurrentBlockHash, ledger.WithTargets(orgTestPeer0.(fab.Peer), orgTestPeer1.(fab.Peer)), ledger.WithMinTargets(2))
	if err != nil {
		log.Fatalf("QueryBlockByHash return error: %v", err)
	}
	log.Printf("Query Block by Hash, retrieve current block by hash:%#v\n", block)
	log.Printf("block CurrentHash:%#v, PreviousHash: %#v, Metadata:%#v\n", string(block.Header.GetDataHash()), string(block.GetHeader().PreviousHash), block.Metadata.String())

	// Org2 user moves funds on org2 peer
	response, err = chClientOrg2User.Execute(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: ExampleCCTxArgs()}, channel.WithProposalProcessor(orgTestPeer1))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	// a-1 b+1
	// Assert that funds have changed value on org1 peer
	verifyValue(chClientOrg1User, initial+1)

	// Get latest block chain info
	ledgerInfoAfter, err := ledgerClient.QueryInfo(ledger.WithTargets(orgTestPeer0.(fab.Peer), orgTestPeer1.(fab.Peer)), ledger.WithMinTargets(2))
	if err != nil {
		log.Fatalf("ledgetInfoAfter QueryInfo return error: %v", err)
	}

	if ledgerInfoAfter.BCI.Height-ledgerInfoBefore.BCI.Height <= 0 {
		log.Fatalf("Block size did not increase after transaction")
	}

	// Test Query Block by Hash - retrieve current block by number
	block, err = ledgerClient.QueryBlock(int(ledgerInfoAfter.BCI.Height)-1, ledger.WithTargets(orgTestPeer0.(fab.Peer), orgTestPeer1.(fab.Peer)), ledger.WithMinTargets(2))
	if err != nil {
		log.Fatalf("QueryBlock return error: %v", err)
	}
	if block == nil {
		log.Fatalf("Block info not available")
	}
	log.Printf("block info : %#v\n", block)

	// Get transaction info
	transactionInfo, err := ledgerClient.QueryTransaction(response.TransactionID, ledger.WithTargets(orgTestPeer0.(fab.Peer), orgTestPeer1.(fab.Peer)), ledger.WithMinTargets(2))
	if err != nil {
		log.Fatalf("QueryTransaction return error: %v", err)
	}
	if transactionInfo.TransactionEnvelope == nil {
		log.Fatalf("Transaction info missing")
	}
	// Start chaincode upgrade process (install and instantiate new version of exampleCC)
	// 链码升级
	installCCReq = resmgmt.InstallCCRequest{Name: "exampleCC", Path: "github.com/example_cc", Version: "1", Package: ccPkg}

	// Install example cc version '1' to Org1 peers
	response1, err = org1ResMgmt.InstallCC(installCCReq)
	if err != nil {
		log.Fatalf("Failed org1 install v2 chaincode:%v", err)
	}
	log.Printf("org1 install v2 chaincode response: %v#\n", response1)

	// Install example cc version '1' to Org2 peers
	response2, err = org2ResMgmt.InstallCC(installCCReq)
	if err != nil {
		log.Fatalf("Failed org2 install v2 chaincode:%v", err)
	}
	log.Printf("org2 install v2 chaincode response: %#v\n", response2)

	// New chaincode policy (both orgs have to approve)
	org1Andorg2Policy, err := cauthdsl.FromString("AND ('Org1MSP.member','Org2MSP.member')")
	if err != nil {
		log.Fatalf("Failed new chaincode policy: %v", err)
	}
	log.Printf("New chaincode policy result:%#v\n", org1Andorg2Policy)

	// Org1 resource manager will instantiate 'example_cc' version 1 on 'orgchannel'
	if err = org1ResMgmt.UpgradeCC(channelID, resmgmt.UpgradeCCRequest{Name: "exampleCC", Path: "github.com/example_cc", Version: "1", Args: ExampleCCUpgradeArgs(), Policy: org1Andorg2Policy}); err != nil {
		log.Fatalf("Failed upgradecc v2 instantiate:%v", err)
	}

	// Org2 user moves funds on org2 peer (cc policy fails since both Org1 and Org2 peers should participate)
	// response, err = chClientOrg2User.Execute(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: ExampleCCTxArgs()}, channel.WithProposalProcessor(orgTestPeer1))
	// if err != nil {
	// 	log.Fatalf("Should have failed to move funds due to cc policy: %v", err)
	// }

	// Org2 user moves funds (cc policy ok since we have provided peers for both Orgs)
	response, err = chClientOrg2User.Execute(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: ExampleCCTxArgs()}, channel.WithProposalProcessor(orgTestPeer0, orgTestPeer1))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	// Assert that funds have changed value on org1 peer
	beforeTxValue, _ := strconv.Atoi(ExampleCCUpgradeB)
	expectedValue := beforeTxValue + 1
	verifyValue(chClientOrg1User, expectedValue)
	return expectedValue
}

func testWithOrg2(expectedValue int) int {
	//指定将由动态选择服务使用的用户（检索chanincode策略信息）
	//此用户必须具有查询链接代码数据的lscc的权限
	mychannelUser := selection.ChannelUser{ChannelID: channelID, UserName: "User1", OrgName: "Org1"}

	//Create SDK setup for channel client with dynamic selection
	sdk, err := fabsdk.New(config.FromFile(configFile), fabsdk.WithServicePkg(&DynamicSelectionProviderFactory{ChannelUsers: []selection.ChannelUser{mychannelUser}}))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()
	// Create new client that will use dynamic selection
	chClientOrg2User, err := sdk.NewClient(fabsdk.WithUser("User1"), fabsdk.WithOrg(org2)).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to create new channel client for Org2 user: %s", err)
	}
	// Org2 user moves funds (dynamic selection will inspect chaincode policy to determine endorsers)
	_, err = chClientOrg2User.Execute(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: ExampleCCTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	expectedValue++
	return expectedValue
}

func verifyWithOrg1(sdk *fabsdk.FabricSDK, expectedValue int) {
	// Org1 user connects to 'orgchannel'
	chClientOrg1User, err := sdk.NewClient(fabsdk.WithUser("User1"), fabsdk.WithOrg(org1)).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to create new channel client for Org1 user: %s", err)
	}
	verifyValue(chClientOrg1User, expectedValue)
}

func loadOrgUser(sdk *fabsdk.FabricSDK, orgName string, userName string) context.Identity {
	session, err := sdk.NewClient(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName)).Session()
	if err != nil {
		log.Fatalf("Session orgName=%s, userName=%s and failed:%s", orgName, userName, err)
	}
	return session
}

func loadOrgPeers(sdk *fabsdk.FabricSDK) {
	org1Peers, err := sdk.Config().PeersConfig(org1)
	if err != nil {
		log.Fatalf("org1 peers failed: %s", err)
	}
	log.Printf("org1Peers data: %#v\n", org1Peers)

	org2Peers, err := sdk.Config().PeersConfig(org2)
	if err != nil {
		log.Fatalf("org2 peers failed: %s", err)
	}
	log.Printf("org2Peers data: %#v\n", org2Peers)

	orgTestPeer0, err = peer.New(sdk.Config(), peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: org1Peers[0]}))
	if err != nil {
		log.Fatal(err)
	}
	orgTestPeer1, err = peer.New(sdk.Config(), peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: org2Peers[0]}))
	if err != nil {
		log.Fatal(err)
	}
}

func verifyValue(chClient *channel.Client, expected int) {
	// Assert that funds have changed value on org1 peer
	var valueInt int
	for i := 0; i < pollRetries; i++ {
		// Query final value on org1 peer
		response, err := chClient.Query(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: ExampleCCQueryArgs()}, channel.WithProposalProcessor(orgTestPeer0))
		if err != nil {
			log.Fatalf("Failed to query funds after transaction: %s", err)
		}
		// If value has not propogated sleep with exponential backoff
		valueInt, _ = strconv.Atoi(string(response.Payload))
		if expected != valueInt {
			backoffFactor := math.Pow(2, float64(i))
			time.Sleep(time.Millisecond * 50 * time.Duration(backoffFactor))
		} else {
			break
		}
	}
	if expected != valueInt {
		log.Fatalf("Org2 'move funds' transaction result was not propagated to Org1. Expected %d, got: %d",
			(expected), valueInt)
	}
}

// DynamicSelectionProviderFactory is configured with dynamic (endorser) selection provider
type DynamicSelectionProviderFactory struct {
	defsvc.ProviderFactory
	ChannelUsers []selection.ChannelUser
}

// CreateSelectionProvider returns a new implementation of dynamic selection provider
func (f *DynamicSelectionProviderFactory) CreateSelectionProvider(config core.Config) (fab.SelectionProvider, error) {
	return selection.New(config, f.ChannelUsers, nil)
}
