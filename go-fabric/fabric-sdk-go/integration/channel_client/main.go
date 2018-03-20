package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/provider/fabpvdr"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

var (
	org1Name          = "Org1"
	chainCodePath     = "../chaincode"
	configFile        = "../../fixtures/config/config_test.yaml"
	channelID         = "mychannel"
	channelConfigPath = "../../fixtures/fabric/v1.0.0/channel/mychannel.tx"
)

//Package fabsdk支持客户端使用Hyperledger Fabric网络。
func main() {
	//func FromFile(name string, opts ...Option) core.ConfigProvider
	//func New(cp core.ConfigProvider, opts ...Option) (*FabricSDK, error)
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatal("SDK init failed:", err)
		return
	}
	log.Printf("*FabricSDK: %#v\n", sdk)
	//func (sdk *FabricSDK) NewClient(identityOpt IdentityOption, opts ...ContextOption) *ClientContext
	// type ClientContext struct {
	// 	provider clientProvider
	// }
	//type clientProvider func() (*clientContext, error)
	// type clientContext struct {
	// 	opts          *contextOptions
	// 	identity      context.Identity
	// 	providers     providers #interface
	// 	clientFactory api.SessionClientFactory
	// }
	//func WithIdentity(identity context.Identity) IdentityOption
	//func WithUser(name string) IdentityOption
	client := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg(org1Name))
	//func (c *ClientContext) Session() (context.Session, error)
	// // Identity supplies the serialized identity and key reference.
	// type Identity interface {
	// 	MspID() string
	// 	SerializedIdentity() ([]byte, error)
	// 	PrivateKey() core.Key
	// }
	//
	// // Session primarily represents the session and identity context
	// type Session interface {
	// 	Identity
	// }
	// // Client supplies the configuration and signing identity to client objects.
	// type Client interface {
	// 	core.Providers
	// 	Identity
	// }
	//
	// // Providers represents the SDK configured providers context.
	// type Providers interface {
	// 	core.Providers
	// 	fab.Providers
	// }
	session, err := client.Session()
	if err != nil {
		log.Fatal("Client session failed:", err)
	}
	log.Printf("context.Session: %#v\n", session)
	log.Printf("MspID %s\n", session.MspID())
	// CreateResourceClient返回一个为SDK当前实例初始化的新客户端。
	// func (sdk *FabricSDK) FabricProvider() fab.InfraProvider
	// type Providers interface package fab   Providers represents the SDK configured service providers context.
	// type Providers interface package core  Providers represents the SDK configured core providers context.
	// CreateResourceClient returns a new client initialized for the current instance of the SDK.
	// func (f *FabricProvider) CreateResourceClient(ic fab.IdentityContext) (api.Resource, error)
	// type IdentityContext interface {
	// 	MspID() string
	// 	SerializedIdentity() ([]byte, error)
	// 	PrivateKey() core.Key
	// }
	// type Resource interface {
	// 	CreateChannel(request CreateChannelRequest) (fab.TransactionID, error)
	// 	InstallChaincode(request InstallChaincodeRequest) ([]*fab.TransactionProposalResponse, fab.TransactionID, error)
	// 	QueryInstalledChaincodes(peer fab.ProposalProcessor) (*pb.ChaincodeQueryResponse, error)
	// 	QueryChannels(peer fab.ProposalProcessor) (*pb.ChannelQueryResponse, error)
	// 	GenesisBlockFromOrderer(channelName string, orderer fab.Orderer) (*common.Block, error)
	// 	LastConfigFromOrderer(channelName string, orderer fab.Orderer) (*common.ConfigEnvelope, error)
	// 	JoinChannel(request JoinChannelRequest) error
	// 	SignChannelConfig(config []byte, signer context.Identity) (*common.ConfigSignature, error)
	// }
	rc, err := sdk.FabricProvider().(*fabpvdr.FabricProvider).CreateResourceClient(session)
	if err != nil {
		log.Fatal("NewResourceClient failed:", err)
	}
	log.Println(rc)
	log.Printf("\r\n")
	targets, err := getOrgTargets(sdk.Config(), org1Name)
	if err != nil {
		log.Fatal("loading target peers from config failed:", err)
	}
	log.Printf("getOrgTargets return targets : %#v\n", targets)
	for _, target := range targets {
		log.Printf("target info: %#v\n", target)
	}
	// Create channel for tests
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfig: channelConfigPath, SigningIdentity: session}
	log.Printf("request params: %#v\n", req)
	log.Println("InitializeChannel start")
	InitializeChannel(sdk, org1Name, req, targets)
	log.Println("InitializeChannel end")
	// Create the channel transactor
	chService, err := client.ChannelService(channelID)
	if err != nil {
		log.Fatal("channel service creation failed: ", err)
	}
	transactor, err := chService.Transactor()
	if err != nil {
		log.Fatal("transactor client creation failed: ", err)
	}
	log.Println(transactor)
	eventHub, err := chService.EventHub()
	if err != nil {
		log.Fatal("eventhub client creation failed:", err)
	}
	if err = eventHub.Connect(); err != nil {
		log.Fatal("eventHub connect failed:", err)
	}
	// test
	chainCodeID := GenerateRandomID()
	if err = InstallAndInstantiateExampleCC(sdk, fabsdk.WithUser("Admin"), org1Name, chainCodeID); err != nil {
		log.Fatal(err)
	}
	actionChainCode(chainCodeID)
}

// Create SDK
func actionChainCode(chainCodeID string) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()

	chClient, err := sdk.NewClient(fabsdk.WithUser("User1")).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to create new channel client: %s", err)
	}
	fmt.Println(chClient)
	// Synchronous query
	testQuery("200", chainCodeID, chClient)

	transientData := "some data"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte(transientData)

	// Synchronous transaction
	response, err := chClient.Execute(
		channel.Request{
			ChaincodeID:  chainCodeID,
			Fcn:          "invoke",
			Args:         ExampleCCTxArgs(),
			TransientMap: transientDataMap,
		})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	// The example CC should return the transient data as a response
	if string(response.Payload) != transientData {
		log.Fatalf("Expecting response [%s] but got [%v]", transientData, response)
	}

	//Verify transaction using query
	log.Println("ChainCodeID:", chainCodeID)
	testQueryWithOpts("201", chainCodeID, chClient)

	// transaction
	testTransaction(chainCodeID, chClient)

	// Verify transaction
	testQuery("202", chainCodeID, chClient)

	// Verify that filter error and commit error did not modify value
	//testQuery("202", chainCodeID, chClient)

	// Test register and receive chaincode event
	testChaincodeEvent(chainCodeID, chClient)

	// Verify transaction with chain code event completed
	testQuery("203", chainCodeID, chClient)

	// Test invocation of custom handler
	testInvokeHandler(chainCodeID, chClient)

	// Test receive event using separate client
	listener, err := sdk.NewClient(fabsdk.WithUser("User1")).Channel(channelID)
	if err != nil {
		log.Fatalf("Failed to create new channel client: %s", err)
	}
	defer listener.Close()
	testChaincodeEventListener(chainCodeID, chClient, listener)

	//Release channel client resources
	if err = chClient.Close(); err != nil {
		log.Fatalf("Failed to close channel client: %v", err)
	}
	log.Println("end success!")
}

func testQuery(expected string, ccID string, chClient *channel.Client) {
	response, err := chClient.Query(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCQueryArgs()})
	if err != nil {
		log.Fatalf("Failed to invoke example cc: %s", err)
	}
	fmt.Println(response)

	if string(response.Payload) != expected {
		log.Fatalf("Expecting %s, got %s", expected, response.Payload)
	}
}

func testQueryWithOpts(expected string, ccID string, chClient *channel.Client) {
	response, err := chClient.Query(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCQueryArgs()})
	if err != nil {
		log.Fatalf("Query returned error: %s", err)
	}
	if string(response.Payload) != expected {
		log.Fatalf("Expecting %s, got %s", expected, response.Payload)
	}
}

func testTransaction(ccID string, chClient *channel.Client) {
	response, err := chClient.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}
	if response.TxValidationCode != pb.TxValidationCode_VALID {
		log.Fatalf("Expecting TxValidationCode to be TxValidationCode_VALID but received: %s", response.TxValidationCode)
	}
}

func testChaincodeEvent(ccID string, chClient *channel.Client) {
	eventID := "test([a-zA-Z]+)"
	//注册chaincode事件（当事件完成时传入接收事件细节的通道）
	notifier := make(chan *channel.CCEvent)
	rce, err := chClient.RegisterChaincodeEvent(notifier, ccID, eventID)
	if err != nil {
		log.Fatalf("Failed to register cc event: %s", err)
	}

	response, err := chClient.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}

	select {
	case ccEvent := <-notifier:
		log.Printf("Received cc event: %s", ccEvent)
		if ccEvent.TxID != string(response.TransactionID) {
			log.Fatalf("CCEvent(%s) and Execute(%s) transaction IDs don't match", ccEvent.TxID, string(response.TransactionID))
		}
	case <-time.After(time.Second * 20):
		log.Fatalf("Did NOT receive CC for eventId(%s)\n", eventID)
	}

	// UnregisterChaincodeEvent 删除链式代码事件注册
	if err = chClient.UnregisterChaincodeEvent(rce); err != nil {
		log.Fatalf("Unregister cc event failed: %s", err)
	}
	log.Println("Event suceess!")
}

func testInvokeHandler(ccID string, chClient *channel.Client) {
	//在提交之前和之后插入自定义处理程序。确保处理程序正在通过写出一些数据并与响应进行比较来调用。
	var txID string
	var endorser string
	txValidationCode := pb.TxValidationCode(-1)
	response, err := chClient.InvokeHandler(
		invoke.NewProposalProcessorHandler(
			invoke.NewEndorsementHandler(
				invoke.NewEndorsementValidationHandler(
					&testHandler{
						txID:     &txID,
						endorser: &endorser,
						next: invoke.NewCommitHandler(
							&testHandler{
								txValidationCode: &txValidationCode,
							},
						),
					},
				),
			),
		),
		channel.Request{
			ChaincodeID: ccID,
			Fcn:         "invoke",
			Args:        ExampleCCTxArgs(),
		},
		channel.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to invoke example cc asynchronously: %s", err)
	}
	log.Println(response)
	if len(response.Responses) == 0 {
		log.Fatalf("Expecting more than one endorsement responses but got none")
	}
	if txID != string(response.TransactionID) {
		log.Fatalf("Expecting TxID [%s] but got [%s]", string(response.TransactionID), txID)
	}
	if endorser != response.Responses[0].Endorser {
		log.Fatalf("Expecting endorser [%s] but got [%s]", response.Responses[0].Endorser, endorser)
	}
	if txValidationCode != response.TxValidationCode {
		log.Fatalf("Expecting TxValidationCode [%s] but got [%s]", response.TxValidationCode, txValidationCode)
	}
}

type testHandler struct {
	txID             *string
	endorser         *string
	txValidationCode *pb.TxValidationCode
	next             invoke.Handler
}

func (h *testHandler) Handle(requestContext *invoke.RequestContext, clientContext *invoke.ClientContext) {
	if h.txID != nil {
		*h.txID = string(requestContext.Response.TransactionID)
		log.Printf("Custom handler writing TxID [%s]", *h.txID)
	}
	if h.endorser != nil && len(requestContext.Response.Responses) > 0 {
		*h.endorser = requestContext.Response.Responses[0].Endorser
		log.Printf("Custom handler writing Endorser [%s]", *h.endorser)
	}
	if h.txValidationCode != nil {
		*h.txValidationCode = requestContext.Response.TxValidationCode
		log.Printf("Custom handler writing TxValidationCode [%s]", *h.txValidationCode)
	}
	if h.next != nil {
		log.Printf("Custom handler invoking next handler")
		h.next.Handle(requestContext, clientContext)
		log.Printf("Custom handler invoking next handler2")
	}
}

func testChaincodeEventListener(ccID string, chClient *channel.Client, listener *channel.Client) {
	eventID := "test([a-zA-Z]+)"

	// Register chaincode event (pass in channel which receives event details when the event is complete)
	notifier := make(chan *channel.CCEvent)
	rce, err := listener.RegisterChaincodeEvent(notifier, ccID, eventID)
	if err != nil {
		log.Fatalf("Failed to register cc event: %s", err)
	}
	log.Println("registerChainCodeEvent ok")
	response, err := chClient.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}

	select {
	case ccEvent := <-notifier:
		log.Printf("Received cc event: %s", ccEvent)
		if ccEvent.TxID != string(response.TransactionID) {
			log.Fatalf("CCEvent(%s) and Execute(%s) transaction IDs don't match", ccEvent.TxID, string(response.TransactionID))
		}
	case <-time.After(20 * time.Second):
		log.Fatalf("Did NOT receive CC for eventId(%s)\n", eventID)
	}

	// Unregister chain code event using registration handle
	if err = listener.UnregisterChaincodeEvent(rce); err != nil {
		log.Fatalf("Unregister cc event failed: %s", err)
	}
	log.Println("testChaincodeEventListener success!")
}

func getOrgTargets(config core.Config, org string) ([]fab.ProposalProcessor, error) {
	targets := []fab.ProposalProcessor{}
	peerConfig, err := config.PeersConfig(org)
	if err != nil {
		log.Fatal("reading peer config failed:", err)
		return nil, err
	}
	for _, p := range peerConfig {
		target, err := peer.New(config, peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: p}))
		if err != nil {
			log.Fatal("NewPeer failed:", err)
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}
