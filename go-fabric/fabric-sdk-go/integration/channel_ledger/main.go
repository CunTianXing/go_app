package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/params"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
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
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	session, err := sdk.NewClient(fabsdk.WithUser(adminUser), fabsdk.WithOrg(org1Name)).Session()
	if err != nil {
		log.Fatalf("failed getting admin user session for org: %s", err)
	}
	targets, err := CreateProposalProcessors(sdk.Config(), []string{org1Name})
	if err != nil {
		log.Fatalf("creating peers failed: %v", err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfig: channelConfigFile, SigningIdentity: session}
	err = InitializeChannel(sdk, org1Name, req, targets)
	if err != nil {
		log.Fatalf("failed to ensure channel has been initialized: %s", err)
	}
	actionLedger(sdk, targets)
}

func actionLedger(sdk *fabsdk.FabricSDK, targets []fab.ProposalProcessor) {
	defer sdk.Close()
	chaincodeID := GenerateRandomID()
	if err := InstallAndInstantiateExampleCC(sdk, fabsdk.WithUser(adminUser), org1Name, chaincodeID); err != nil {
		log.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
	}
	//Get a ledger Client
	client := sdk.NewClient(fabsdk.WithUser(adminUser), fabsdk.WithOrg(org1Name))
	channelSvc, err := client.ChannelService(channelID)
	if err != nil {
		log.Fatalf("creating channel service failed: %v", err)
	}
	ledger, err := channelSvc.Ledger()
	if err != nil {
		log.Fatalf("creating channel ledger client failed: %v", err)
	}

	// Test Query Info - retrieve values before transaction
	bciBeforeTx, err := ledger.QueryInfo(targets[0:1])
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}

	// Invoke transaction that changes block state
	channel, err := client.Channel(channelID)
	if err != nil {
		log.Fatalf("creating channel failed: %v", err)
	}

	txID, err := changeBlockState(channel, chaincodeID)
	if err != nil {
		log.Fatalf("QueryInfo return error: %s", err)
	}
	log.Println("txID:", txID)
	// Test Query Info - retrieve values after transaction
	bciAfterTx, err := ledger.QueryInfo(targets[0:1])
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}

	// Test Query Info -- verify block size changed after transaction
	if (bciAfterTx[0].BCI.Height - bciBeforeTx[0].BCI.Height) <= 0 {
		log.Fatalf("Block size did not increase after transaction")
	}
	testQueryTransaction(ledger, txID, targets)
	testQueryBlock(ledger, targets)
	testInstantiatedChaincodes(chaincodeID, ledger, targets)
	testQueryConfigBlock(ledger, targets)
}

func changeBlockState(client *channel.Client, chaincodeID string) (fab.TransactionID, error) {
	req := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         "invoke",
		Args:        params.ExampleCCQueryArgs(),
	}
	resp, err := client.Query(req)
	if err != nil {
		return "", err
	}
	value := resp.Payload

	txID, err := moveFundsAndGetTxID(client, chaincodeID)
	if err != nil {
		return "", err
	}

	resp, err = client.Query(req)
	if err != nil {
		return "", err
	}
	valueAfterInvoke := resp.Payload

	// Verify that transaction changed block state
	valueInt, _ := strconv.Atoi(string(value))
	valueInt = valueInt + 1
	valueAfterInvokeInt, _ := strconv.Atoi(string(valueAfterInvoke))
	if valueInt != valueAfterInvokeInt {
		log.Printf("SendTransaction didn't change the QueryValue %s", value)
		return "", err
	}
	return txID, nil
}

func testQueryConfigBlock(ledger fab.ChannelLedger, targets []fab.ProposalProcessor) {
	//检索当前通道配置
	cfgEnvelope, err := ledger.QueryConfigBlock(targets, 1)
	if err != nil {
		log.Fatalf("QueryConfig return error: %v", err)
	}
	if cfgEnvelope.Config == nil {
		log.Fatalf("QueryConfig config data is nil")
	}
	log.Printf("channel config: %#v\n", cfgEnvelope.GetConfig().String())
}

func testInstantiatedChaincodes(ccID string, ledger fab.ChannelLedger, targets []fab.ProposalProcessor) {
	// Test Query Instantiated chaincodes
	chaincodeQueryResponses, err := ledger.QueryInstantiatedChaincodes(targets)
	if err != nil {
		log.Fatalf("QueryInstantiatedChaincodes return error: %v", err)
	}
	found := false
	for _, chaincodeQueryResponse := range chaincodeQueryResponses {
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			log.Printf("**InstantiatedCC: %s", chaincode)
			if chaincode.Name == ccID {
				found = true
			}
		}
	}
	if !found {
		log.Fatalf("QueryInstantiatedChaincodes failed to find instantiated %s chaincode", ccID)
	}
	log.Println("ok")
}

func testQueryBlock(ledger fab.ChannelLedger, targets []fab.ProposalProcessor) {
	//检索当前的区块链信息
	bcis, err := ledger.QueryInfo(targets)
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}
	for i, bci := range bcis {
		//通过哈希测试查询块 - 通过哈希检索当前块
		blocks, err := ledger.QueryBlockByHash(bci.BCI.CurrentBlockHash, targets[i:i+1])
		if err != nil {
			log.Fatalf("QueryBlockByHash return error: %v", err)
		}
		for _, block := range blocks {
			log.Printf("block Number %d\n", block.GetHeader().GetNumber())
		}

		if blocks[0].Data == nil {
			log.Println("QueryBlockByHash block data is nil")
		}
	}
	// Test Query Block - retrieve block by number
	blockss, err := ledger.QueryBlock(1, targets)
	if err != nil {
		log.Fatalf("QueryBlock return error: %v", err)
	}
	for _, b := range blockss {
		if b.Data == nil {
			log.Fatalf("QueryBlock block data is nil")
		}
	}
}

// ProposalProcessor模拟交易提案，以便客户端可以提交结果进行排序。
func testQueryTransaction(ledger fab.ChannelLedger, txID fab.TransactionID, targets []fab.ProposalProcessor) {
	//测试查询交易 - 确认有效的交易已被处理
	processedTransactions, err := ledger.QueryTransaction(txID, targets)
	if err != nil {
		log.Fatalf("QueryTransaction return error: %v\n", err)
	}

	for _, processedTransaction := range processedTransactions {
		if processedTransaction.TransactionEnvelope == nil {
			log.Fatalf("QueryTransaction failed to return transaction envelope")
		}
	}

	//测试查询交易 - 检索不存在的交易
	// _, err = ledger.QueryTransaction("dddddsd", targets)
	// if err == nil {
	// 	log.Fatalf("QueryTransaction non-existing didn't return an error")
	// }

}

func moveFundsAndGetTxID(client *channel.Client, chaincodeID string) (fab.TransactionID, error) {
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")
	req := channel.Request{
		ChaincodeID:  chaincodeID,
		Fcn:          "invoke",
		Args:         params.ExampleCCTxArgs(),
		TransientMap: transientDataMap,
	}
	resp, err := client.Execute(req)
	if err != nil {
		return "", err
	}
	return resp.TransactionID, nil
}

//Initialize...
func Initialize(sdk *fabsdk.FabricSDK) error {
	client := sdk.NewClient(fabsdk.WithUser(adminUser), fabsdk.WithOrg(org1Name))
	session, err := client.Session()
	if err != nil {
		return err
	}
	// resClient, err := sdk.FabricProvider().(*fabpvdr.FabricProvider).CreateResourceClient(session)
	// if err != nil {
	// 	return err
	// }
	targets, err := getOrgTargets(sdk.Config(), org1Name)
	if err != nil {
		log.Fatalf("get peers failed: %s", err)
	}
	//Create channel
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfig: channelConfigFile, SigningIdentity: session}
	err = InitializeChannel(sdk, org1Name, req, targets)
	if err != nil {
		log.Fatalf("create channel and join peer failed: %s", err)
	}

	//Create the channel transactor
	chService, err := client.ChannelService(channelID)
	if err != nil {
		log.Fatalf("create the channel transactor failed:%s", err)
	}
	transactor, err := chService.Transactor()
	if err != nil {
		log.Fatalf("transactor client creation failed:%s", err)
	}
	log.Printf("transactor:%#v\n", transactor)
	eventHub, err := chService.EventHub()
	if err != nil {
		log.Fatalf("eventhub client creation failed: %s", err)
	}
	if err := eventHub.Connect(); err != nil {
		log.Fatalf("eventHub connect failed:%s", err)
	}
	return nil
}

// InitializeChannel ...
func InitializeChannel(sdk *fabsdk.FabricSDK, orgID string, req resmgmt.SaveChannelRequest, targets []fab.ProposalProcessor) error {
	joinTargets, err := FilterTargetsJoinedChannel(sdk, orgID, req.ChannelID, targets)
	if err != nil {
		return err
	}
	if len(joinTargets) != len(targets) {
		_, err := CreateChannel(sdk, req)
		if err != nil {
			return err
		}
		_, err = JoinChannel(sdk, req.ChannelID)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateChannel attempts to save the named channel.
func CreateChannel(sdk *fabsdk.FabricSDK, req resmgmt.SaveChannelRequest) (bool, error) {
	// Channel management client is responsible for managing channels (create/update)
	resMgmtClient, err := sdk.NewClient(fabsdk.WithUser(adminUser), fabsdk.WithOrg(org1Name)).ResourceMgmt()
	if err != nil {
		return false, err
	}
	// Create channel (or update if it already exists)
	if err = resMgmtClient.SaveChannel(req); err != nil {
		return false, nil
	}
	time.Sleep(time.Second * 5)
	return true, nil
}

func JoinChannel(sdk *fabsdk.FabricSDK, name string) (bool, error) {
	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := sdk.NewClient(fabsdk.WithUser(adminUser)).ResourceMgmt()
	if err != nil {
		return false, err
	}
	if err = resMgmtClient.JoinChannel(name); err != nil {
		return false, err
	}
	return true, nil
}

func FilterTargetsJoinedChannel(sdk *fabsdk.FabricSDK, orgID string, channelID string, targets []fab.ProposalProcessor) ([]fab.ProposalProcessor, error) {
	joinedTargets := []fab.ProposalProcessor{}
	client, err := sdk.NewClient(fabsdk.WithUser(adminUser), fabsdk.WithOrg(orgID)).ResourceMgmt()
	if err != nil {
		log.Printf("failed getting admin user session for org:%s", err)
		return nil, err
	}
	for _, peer := range targets {
		// Check if primary peer has joined channel
		alreadJoined, err := HasPeerJoinedChannel(client, peer, channelID)
		if err != nil {
			return nil, err
		}
		if alreadJoined {
			joinedTargets = append(joinedTargets, peer)
		}
	}
	return joinedTargets, nil
}

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

func getOrgTargets(config core.Config, org string) ([]fab.ProposalProcessor, error) {
	targets := []fab.ProposalProcessor{}
	peerConfig, err := config.PeersConfig(org)
	if err != nil {
		log.Printf("reading peer config failed:%s", err)
		return nil, err
	}
	for _, p := range peerConfig {
		target, err := peer.New(config, peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: p}))
		if err != nil {
			log.Printf("NewPeer failed:%s", err)
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func InstallAndInstantiateExampleCC(sdk *fabsdk.FabricSDK, user fabsdk.IdentityOption, orgName string, chainCodeID string) error {
	return InstallAndInstantiateCC(sdk, user, orgName, chainCodeID, "github.com/example_cc", "v0", chainCodePath, params.ExampleCCInitArgs())
}

func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, user fabsdk.IdentityOption, orgName string, ccName, ccPath, ccVersion, goPath string, ccArgs [][]byte) error {
	ccPkg, err := packager.NewCCPackage(ccPath, goPath)
	if err != nil {
		log.Printf("creating chaincode package failed:%s", err)
		return err
	}
	mspID, err := sdk.Config().MspID(orgName)
	if err != nil {
		log.Printf("looking up MSP ID failed:%s", err)
		return err
	}

	resMgmtClient, err := sdk.NewClient(user, fabsdk.WithOrg(orgName)).ResourceMgmt()
	if err != nil {
		log.Printf("Failed to create new resource management client:%s", err)
		return err
	}

	_, err = resMgmtClient.InstallCC(resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg})
	if err != nil {
		log.Printf("install chaincode failed:%s", err)
		return err
	}
	ccPolicy := cauthdsl.SignedByMspMember(mspID)
	return resMgmtClient.InstantiateCC(channelID, resmgmt.InstantiateCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Args: ccArgs, Policy: ccPolicy})
}

// CreateProposalProcessors initializes target peers based on config
func CreateProposalProcessors(config core.Config, orgs []string) ([]fab.ProposalProcessor, error) {
	peers := []fab.ProposalProcessor{}
	for _, org := range orgs {
		peerConfig, err := config.PeersConfig(org)
		if err != nil {
			log.Printf("reading peer config failed:%s", err)
			return nil, err
		}
		for _, p := range peerConfig {
			endorser, err := peer.New(config, peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: p}))
			if err != nil {
				log.Printf("new peer failed: %s", err)
				return nil, err
			}
			peers = append(peers, endorser)
		}
	}
	return peers, nil
}

//GenerateRandomID generates random ID
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
