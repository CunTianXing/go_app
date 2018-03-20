package comm

import (
	"log"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/txn"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
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

// InitializeChannel ...
func InitializeChannel(sdk *fabsdk.FabricSDK, orgID string, req resmgmt.SaveChannelRequest, targets []fab.ProposalProcessor) error {
	joinedTargets, err := FilterTargetsJoinedChannel(sdk, orgID, req.ChannelID, targets)
	if err != nil {
		log.Printf("checking for joined targets failed:%s\n", err)
		return err
	}

	if len(joinedTargets) != len(targets) {
		_, err := CreateChannel(sdk, req)
		if err != nil {
			log.Printf("create channel failed:%s\n", err)
			return err
		}

		_, err = JoinChannel(sdk, req.ChannelID)
		if err != nil {
			log.Printf("join channel failed: %s\n", err)
			return err
		}
	}
	return nil
}

// InstallAndInstantiateCC install and instantiate using resource management client
func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, user fabsdk.IdentityOption, orgName string, channelID, ccName, ccPath, ccVersion, goPath string, ccArgs [][]byte) error {

	ccPkg, err := packager.NewCCPackage(ccPath, goPath)
	if err != nil {
		log.Printf("creating chaincode package failed:%s\n", err)
		return err
	}

	mspID, err := sdk.Config().MspID(orgName)
	if err != nil {
		log.Printf("looking up MSP ID failed:%s\n", err)
		return err
	}

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := sdk.NewClient(user, fabsdk.WithOrg(orgName)).ResourceMgmt()
	if err != nil {
		log.Printf("Failed to create new resource management client:%s\n", err)
		return err
	}

	_, err = resMgmtClient.InstallCC(resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg})
	if err != nil {
		return err
	}

	ccPolicy := cauthdsl.SignedByMspMember(mspID)

	return resMgmtClient.InstantiateCC(channelID, resmgmt.InstantiateCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Args: ccArgs, Policy: ccPolicy})
}

// FilterTargetsJoinedChannel filters targets to those that have joined the named channel.
func FilterTargetsJoinedChannel(sdk *fabsdk.FabricSDK, orgID string, channelID string, targets []fab.ProposalProcessor) ([]fab.ProposalProcessor, error) {
	joinedTargets := []fab.ProposalProcessor{}
	rc, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgID)).ResourceMgmt()
	if err != nil {
		log.Printf("failed getting admin user session for org: %s\n", err)
		return nil, err
	}

	for _, target := range targets {
		// Check if primary peer has joined channel
		alreadyJoined, err := HasPeerJoinedChannel(rc, target, channelID)
		if err != nil {
			log.Printf("failed while checking if primary peer has already joined channel:%s\n", err)
			return nil, err
		}
		if alreadyJoined {
			joinedTargets = append(joinedTargets, target)
		}
	}
	return joinedTargets, nil
}

// CreateChannel attempts to save the named channel.
func CreateChannel(sdk *fabsdk.FabricSDK, req resmgmt.SaveChannelRequest) (bool, error) {
	// Channel management client is responsible for managing channels (create/update)
	resMgmtClient, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg("ordererorg")).ResourceMgmt()
	if err != nil {
		log.Printf("Failed to create new channel management client:%s\n", err)
		return false, err
	}

	// Create channel (or update if it already exists)
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
		log.Printf("Failed to create new resource management client:%s\n", err)
		return false, err
	}

	if err = resMgmtClient.JoinChannel(name); err != nil {
		return false, nil
	}
	return true, nil
}

// HasPeerJoinedChannel checks whether the peer has already joined the channel.
// It returns true if it has, false otherwise, or an error
func HasPeerJoinedChannel(client *resmgmt.Client, peer fab.ProposalProcessor, channel string) (bool, error) {
	foundChannel := false
	response, err := client.QueryChannels(peer)
	if err != nil {
		log.Printf("failed to query channel for peer:%s\n", err)
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
		log.Printf("Reading peer config failed: %s\n", err)
		return nil, err
	}
	for _, p := range peerConfig {
		target, err := peer.New(config, peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: p}))
		if err != nil {
			log.Printf("NewPeer failed: %s\n", err)
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

// CreateAndSendTransactionProposal ...
// ProposalProcessor模拟交易提议，以便客户端可以提交结果进行排序。
// TransactionProposal包含一个marashalled事务提议。
// ChaincodeInvokeRequest包含发送交易提案的参数。
//ProposalSender提供了创建和发送交易提议的功能....
func CreateAndSendTransactionProposal(transactor fab.ProposalSender, chainCodeID string, fcn string, args [][]byte, targets []fab.ProposalProcessor, transientData map[string][]byte) ([]*fab.TransactionProposalResponse, *fab.TransactionProposal, error) {
	propReq := fab.ChaincodeInvokeRequest{
		Fcn:          fcn,
		Args:         args,
		TransientMap: transientData,
		ChaincodeID:  chainCodeID,
	}
	//CreateTransactionHeader根据当前上下文创建一个Transaction Header。
	//TransactionHeader提供了一个事务元数据的句柄。
	txh, err := transactor.CreateTransactionHeader()
	if err != nil {
		log.Fatalf("creating transaction header failed:%s", err)
		return nil, nil, err
	}
	//CreateChaincodeInvokeProposal为交易创建提案。
	// TransactionProposal包含一个marashalled事务提议。
	tp, err := txn.CreateChaincodeInvokeProposal(txh, propReq)
	if err != nil {
		log.Fatalf("creating transaction proposal failed:%s", err)
		return nil, nil, err
	}
	//SendTransactionProposal发送一个TransactionProposal给目标节点。
	//交易提案响应表示交易提案处理的结果。
	tpr, err := transactor.SendTransactionProposal(tp, targets)
	return tpr, tp, err
}

//RegisterTxEvent ...
//给定eventhub上的RegisterTxEvent寄存器用于给定事务返回一个布尔通道，
//当事件完成时接收true并且出现错误的错误通道
func RegisterTxEvent(txID fab.TransactionID, eventHub fab.EventHub) (chan bool, chan error) {
	done := make(chan bool)
	fail := make(chan error)
	eventHub.RegisterTxEvent(txID, func(txId fab.TransactionID, errorCode pb.TxValidationCode, err error) {
		if err != nil {
			log.Printf("Received error event for txid(%s)", txId)
			fail <- err
		} else {
			log.Printf("Received success event for txid(%s)", txId)
			done <- true
		}
	})
	return done, fail
}

// CreateAndSendTransaction ...
func CreateAndSendTransaction(transactor fab.Sender, proposal *fab.TransactionProposal, resps []*fab.TransactionProposalResponse) (*fab.TransactionResponse, error) {
	txRequest := fab.TransactionRequest{
		Proposal:          proposal,
		ProposalResponses: resps,
	}
	tx, err := transactor.CreateTransaction(txRequest)
	if err != nil {
		log.Printf("CreateTransaction failed:%s\n", err)
		return nil, err
	}

	transactionResponse, err := transactor.SendTransaction(tx)
	if err != nil {
		log.Printf("SendTransaction failed:%s\n", err)
		return nil, err
	}
	return transactionResponse, nil
}
