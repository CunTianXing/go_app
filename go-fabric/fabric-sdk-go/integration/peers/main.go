package main

import (
	"log"
	"strconv"
	"time"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/params"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

const (
	org1 = "Org1"
	org2 = "Org2"
)

var (
	configFile        = "../../fixtures/config/config_test.yaml"
	channelID         = "orgchannel"
	channelConfigFile = "../../fixtures/fabric/v1.0.0/channel/orgchannel.tx"
	chainCodePath     = "../chaincode"
)

//
func main() {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()
	actionChannel(sdk)
}

func actionChannel(sdk *fabsdk.FabricSDK) {
	// 通道管理客户端负责管理管理（创建/更新通道）
	chMgmtClient, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg("ordererorg")).ResourceMgmt()
	if err != nil {
		log.Fatalf("Failed to new channel management client: %s", err)
	}
	// 创建通道或者更新已经存在的通道
	org1AdminUser := loadOrgUser(sdk, org1, "Admin")
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfig: channelConfigFile, SigningIdentity: org1AdminUser}
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
	// GetPeers by org1
	peers1, err := getPeers(sdk, org1)
	if err != nil {
		log.Fatalf("get peers failed :%s", err)
	}
	// Org1 peers join channel
	if err = org1ResMgmt.JoinChannel(channelID, resmgmt.WithTargets(peers1...)); err != nil {
		log.Fatalf("Org1 all peers failed to JoinChannel: %s", err)
	}
	log.Printf("Org1 all peers add channel ok")
	// Org2资源管理客户端
	org2ResMgmt, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg(org2)).ResourceMgmt()
	if err != nil {
		log.Fatalf("Failed to create new org2 resource management client: %s", err)
	}
	// GetPeers by org2
	peers2, err := getPeers(sdk, org2)
	if err != nil {
		log.Fatalf("get peers failed :%s", err)
	}
	// Org2 加入通道 channelID
	if err = org2ResMgmt.JoinChannel(channelID, resmgmt.WithTargets(peers2...)); err != nil {
		log.Fatalf("Org2 all peers failed to JoinChannel: %s", err)
	}
	log.Printf("Org2 all peers add channel ok")
	// Create chaincode package
	ccPkg, err := packager.NewCCPackage("github.com/example_cc", chainCodePath)
	if err != nil {
		log.Fatalf("Failed to create chaincode package: %s", err)
	}
	installCCReq := resmgmt.InstallCCRequest{Name: "exampleCC", Path: "github.com/example_cc", Version: "0", Package: ccPkg}
	// Install example cc to Org1 peers
	response1, err := org1ResMgmt.InstallCC(installCCReq, resmgmt.WithTargets(peers1...))
	if err != nil {
		log.Fatalf("Failed to install example cc to Org1 peers: %s", err)
	}
	log.Printf("Install example cc to Org1 peers response: %#v\n", response1)
	// Install example cc to Org2 peers
	response2, err := org2ResMgmt.InstallCC(installCCReq, resmgmt.WithTargets(peers2...))
	if err != nil {
		log.Fatalf("Failed to install example cc to Org2 peers: %s", err)
	}
	log.Printf("Install example cc to Org2 peers response: %#v\n", response2)
	log.Println("================Install chaincode ok====================")
	log.Println("Set up chaincode policy to 'any of two msps'")
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})
	// Org1 资源管理者将实例化 'example_cc' 在 'orgchannel'
	log.Printf("org1 instantiate peer is %s", peers1[0])
	if err = org1ResMgmt.InstantiateCC(channelID, resmgmt.InstantiateCCRequest{Name: "exampleCC", Path: "github.com/example_cc", Version: "0", Args: params.ExampleCCInitArgs(), Policy: ccPolicy}, resmgmt.WithTarget(peers1[0])); err != nil {
		log.Fatalf("Failed to org1 instantiate example_cc: %s", err)
	}
	log.Println("org1 instantiate success!")
	// Org2 user connects to orgchannel
	chClientOrg2User, err := sdk.NewClient(fabsdk.WithUser("User1"), fabsdk.WithOrg(org2)).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to org2 user connect to orgchannel: %s", err)
	}
	log.Printf("Org2 user connects to '%s' success!", channelID)
	// Org2用户查询Org1两个对等节点的初始值
	log.Printf("org2 query peer is %s", peers2[0])
	response, err := chClientOrg2User.Query(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: params.ExampleCCQueryArgs()}, channel.WithProposalProcessor(peers2[0]))
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}
	log.Printf("org2 user query initial value on both peers: %#v\n", response)
	initial, _ := strconv.Atoi(string(response.Payload))
	log.Println("Payload:", initial)
	log.Printf("Org2 user move funds on org2 peer %s\n", peers2[1])
	response, err = chClientOrg2User.Execute(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: params.ExampleCCTxArgs()}, channel.WithProposalProcessor(peers2[1]))
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	log.Printf("Org2 user move funds :%#v\n", response)
	//Org1 user connects to orgchannel
	chClientOrg1User, err := sdk.NewClient(fabsdk.WithUser("User1"), fabsdk.WithOrg(org1)).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to org1 user connect to orgchannel: %s", err)
	}
	log.Printf("org1 user to connects to '%s' success!\n", channelID)
	log.Printf("org1 user query on org1 peer %s\n", peers1[1])
	response, err = chClientOrg1User.Query(channel.Request{ChaincodeID: "exampleCC", Fcn: "invoke", Args: params.ExampleCCQueryArgs()}, channel.WithProposalProcessor(peers1[1]))
	if err != nil {
		log.Fatalf("Failed to query funds:%s", err)
	}
	log.Printf("org1 user query initial value on both peers: %#v\n", response)
	value, _ := strconv.Atoi(string(response.Payload))
	log.Println("Payload ", value)
}

func getPeers(sdk *fabsdk.FabricSDK, orgName string) ([]fab.Peer, error) {
	peers := []fab.Peer{}
	orgPeers, err := sdk.Config().PeersConfig(orgName)
	if err != nil {
		log.Printf("get peersConfig by orgName failed: %s", err)
		return nil, err
	}
	for _, peerConfig := range orgPeers {
		peer, err := peer.New(sdk.Config(), peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: peerConfig}))
		if err != nil {
			log.Printf("get peer by peerConfig failed: %s", err)
			return nil, err
		}
		peers = append(peers, peer)
	}
	return peers, nil
}

func loadOrgUser(sdk *fabsdk.FabricSDK, orgName string, userName string) context.Identity {
	session, err := sdk.NewClient(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName)).Session()
	if err != nil {
		log.Fatalf("Session orgName=%s, userName=%s and failed:%s", orgName, userName, err)
	}
	return session
}
