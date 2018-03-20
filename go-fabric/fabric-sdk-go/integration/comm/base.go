package comm

import (
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/provider/fabpvdr"
)

// Common ...
type Common struct {
	SDK             *fabsdk.FabricSDK
	Identity        context.Identity
	Client          api.Resource
	Transactor      fab.Transactor
	Targets         []fab.ProposalProcessor
	EventHub        fab.EventHub
	ConnectEventHub bool
	ConfigFile      string
	OrgID           string
	AdminUser       string
	ChannelID       string
	GoPath          string
	Initialized     bool
	ChannelConfig   string
}

// Initialize
func (c *Common) Initialize() error {
	// Create SDK setup
	sdk, err := fabsdk.New(config.FromFile(c.ConfigFile))
	if err != nil {
		log.Printf("SDK init failed: %s\n", err)
		return err
	}
	c.SDK = sdk
	client := sdk.NewClient(fabsdk.WithUser(c.AdminUser), fabsdk.WithOrg(c.OrgID))
	session, err := client.Session()
	if err != nil {
		log.Printf("Failed getting admin user session for org:%s\n", err)
		return err
	}
	c.Identity = session
	rc, err := sdk.FabricProvider().(*fabpvdr.FabricProvider).CreateResourceClient(c.Identity)
	if err != nil {
		log.Printf("NewResourceClient failed: %s\n", err)
		return err
	}
	c.Client = rc
	targets, err := getOrgTargets(c.SDK.Config(), c.OrgID)
	if err != nil {
		log.Printf("Loading target peers from config failed: %s\n", err)
		return err
	}
	c.Targets = targets
	req := resmgmt.SaveChannelRequest{ChannelID: c.ChannelID, ChannelConfig: c.ChannelConfig, SigningIdentity: c.Identity}
	InitializeChannel(c.SDK, c.OrgID, req, c.Targets)
	// Create the channel transactor
	chService, err := client.ChannelService(c.ChannelID)
	if err != nil {
		log.Printf("channel service creation failed:%s\n", err)
		return err
	}
	transactor, err := chService.Transactor()
	if err != nil {
		log.Printf("transactor client creation failed: %s\n", err)
		return err
	}
	c.Transactor = transactor

	eventHub, err := chService.EventHub()
	if err != nil {
		log.Printf("eventhub client creation failed: %s\n", err)
		return err
	}
	if c.ConnectEventHub {
		if err := eventHub.Connect(); err != nil {
			log.Printf("eventHub connect failed: %s\n", err)
			return err
		}
	}
	c.EventHub = eventHub
	c.Initialized = true
	return nil
}
