package main

import (
	"log"
	"strings"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/chconfig"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defcore"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/provider/fabpvdr"
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
	c := comm.Common{
		ConfigFile:      configFile,
		ChannelID:       channelID,
		ChannelConfig:   channelConfigFile,
		OrgID:           org1Name,
		AdminUser:       adminUser,
		GoPath:          chainCodePath,
		ConnectEventHub: true,
	}
	if err := c.Initialize(); err != nil {
		log.Fatalf(err.Error())
	}
	defer c.SDK.Close()
	TestChannelConfig(c.SDK, c)
	confProvider := config.FromFile(c.ConfigFile)
	//使用orderer的检索通道配置为通道客户端创建SDK设置
	address := "orderer.example.com:7050"
	sdk, err := fabsdk.New(confProvider, fabsdk.WithCorePkg(&ChannelConfigFromOrdererProviderFactory{orderer: setupOrderer(confProvider, address)}))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()
	TestChannelConfig(sdk, c)
}

//TestChannelConfig ....
func TestChannelConfig(sdk *fabsdk.FabricSDK, c comm.Common) {
	// ChannelService返回一个用于与通道交互的客户端API。
	cs, err := sdk.NewClient(fabsdk.WithUser("User1")).ChannelService(c.ChannelID)
	if err != nil {
		log.Fatalf("Failed to create new channel service: %s", err)
	}

	cfg, err := cs.Config()
	if err != nil {
		log.Fatalf("Failed to create new channel config: %s", err)
	}

	response, err := cfg.Query()
	if err != nil {
		log.Fatalf("cfg query err:%s\n", err.Error())
	}

	for _, anchorPeer := range response.AnchorPeers() {
		log.Printf("anchorPeer :%v\n", anchorPeer)
	}
	log.Println(response.AnchorPeers())

	expected := "orderer.example.com:7050"
	found := false
	for _, o := range response.Orderers() {
		log.Printf("orderer: %v\n", o)
		if o == expected {
			found = true
			break
		}
	}

	if !found {
		log.Fatalf("Expected orderer %s, got %s", expected, response.Orderers())
	}
}

func setupOrderer(confProvider core.ConfigProvider, address string) fab.Orderer {
	conf, err := confProvider()
	if err != nil {
		log.Fatalf("confProvider: %s\n", err)
	}

	//Get orderer config by orderer address
	oCfg, err := conf.OrdererConfig(resolveOrdererAddress(address))
	if err != nil {
		log.Fatalf("OrdererConfig err: %s\n", err)
	}

	o, err := orderer.New(conf, orderer.FromOrdererConfig(oCfg))
	if err != nil {
		log.Fatalf("Orderer New err : %s\n", err)
	}
	return o
}

//resolveOrdererAddress ...
func resolveOrdererAddress(ordererAddress string) string {
	s := strings.Split(ordererAddress, ":")
	if len(s) > 1 {
		return s[0]
	}
	return ordererAddress
}

// ChannelConfigFromOrdererProviderFactory is configured to retrieve channel config from orderer
type ChannelConfigFromOrdererProviderFactory struct {
	defcore.ProviderFactory
	orderer fab.Orderer
}

func (f *ChannelConfigFromOrdererProviderFactory) CreateFabricProvider(context core.Providers) (fab.InfraProvider, error) {
	fabProvider := fabpvdr.New(context)
	cfg := CustomFabricProvider{
		FabricProvider:  fabProvider,
		providerContext: context,
		orderer:         f.orderer,
	}
	return &cfg, nil
}

// CustomFabricProvider overrides channel config default implementation
type CustomFabricProvider struct {
	*fabpvdr.FabricProvider
	orderer         fab.Orderer
	providerContext core.Providers
}

// CreateChannelConfig initializes the channel config
func (f *CustomFabricProvider) CreateChannelConfig(ic fab.IdentityContext, channelID string) (fab.ChannelConfig, error) {
	ctx := chconfig.Context{
		Providers: f.providerContext,
		Identity:  ic,
	}
	return chconfig.New(ctx, channelID, chconfig.WithOrderer(f.orderer))
}
