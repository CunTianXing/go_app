package integration

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	pfab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// Initial B values for ExampleCC
const (
	ExampleCCInitB    = "200"
	ExampleCCUpgradeB = "400"
	AdminUser         = "Admin"
)

// ExampleCC query and transaction arguments
var queryArgs = [][]byte{[]byte("query"), []byte("b")}
var txArgs = [][]byte{[]byte("move"), []byte("a"), []byte("b"), []byte("1")}

// ExampleCC init and upgrade args
var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte(ExampleCCInitB)}
var upgradeArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte(ExampleCCUpgradeB)}

//GetInitArgs ...
func GetInitArgs() [][]byte {
	return initArgs
}

//GetQueryArgs ...
func GetQueryArgs() [][]byte {
	return queryArgs
}

// GetTxArgs returns example cc move funds args
func GetTxArgs() [][]byte {
	return txArgs
}

// GetUpgradeArgs ...
func GetUpgradeArgs() [][]byte {
	return upgradeArgs
}

//BaseConfig ....
type BaseConfig struct {
	Identity          msp.Identity
	Targets           []string
	ConfigFile        string
	OrgID             string
	ChannelID         string
	ChannelConfigFile string
}

//Initialize ...
func (b *BaseConfig) Initialize(sdk *fabsdk.FabricSDK) error {
	adminIdentity, err := GetSigningIdentity(sdk, AdminUser, b.OrgID)
	if err != nil {
		log.Printf("failed to get client context %s", err)
		return err
	}
	b.Identity = adminIdentity
	configBackend, err := sdk.Config()
	if err != nil {
		configBackend, err = config.FromFile(b.ConfigFile)()
		if err != nil {
			log.Printf("failed to get config backend from config: %v", err)
			return err
		}
	}
	targets, err := getOrgTargets(configBackend, b.OrgID)
	if err != nil {
		log.Printf("loading target peers from config failed %v\n", err)
		return err
	}
	b.Targets = targets
	r, err := os.Open(b.ChannelConfigFile)
	if err != nil {
		log.Printf("opening channel config file failed: %v\n", err)
		return err
	}
	defer r.Close()
	req := resmgmt.SaveChannelRequest{ChannelID: b.ChannelID, ChannelConfig: r, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	if err = InitializeChannel(sdk, b.OrgID, req, targets); err != nil {
		log.Printf("failed to initialize channel: %v\n", err)
		return err
	}
	return nil
}

func getOrgTargets(configBackend core.ConfigBackend, org string) ([]string, error) {
	endpointConfig, err := fab.ConfigFromBackend(configBackend)
	if err != nil {
		log.Printf("reading config failed %v\n", err)
		return nil, err
	}

	var targets []string

	peerConfig, err := endpointConfig.PeersConfig(org)
	if err != nil {
		log.Printf("reading peer config failed %v\n", err)
		return nil, err
	}

	for _, p := range peerConfig {
		targets = append(targets, p.URL)
	}
	return targets, nil

}

//InitConfig ....
func (b *BaseConfig) InitConfig() core.ConfigProvider {
	return config.FromFile(b.ConfigFile)
}

//InstallAndInstantiateCC ...
func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, user fabsdk.ContextOption, orgName, channelID, ccName, ccPath, ccVersion, goPath string, ccArgs [][]byte) (resmgmt.InstantiateCCResponse, error) {
	ccPkg, err := packager.NewCCPackage(ccPath, goPath)
	if err != nil {
		log.Printf("creating chaincode package failed: %v\n", err)
		return resmgmt.InstantiateCCResponse{}, err
	}

	configBackend, err := sdk.Config()
	if err != nil {
		log.Printf("failed to get config backend:%v\n", err)
		return resmgmt.InstantiateCCResponse{}, err
	}

	endpointConfig, err := fab.ConfigFromBackend(configBackend)
	if err != nil {
		log.Printf("failed to get endpoint config:%v\n", err)
		return resmgmt.InstantiateCCResponse{}, err
	}

	mspID, err := endpointConfig.MSPID(orgName)
	if err != nil {
		log.Printf("looking up MSP ID failed:%v\n", err)
		return resmgmt.InstantiateCCResponse{}, err
	}
	log.Println(mspID)
	clientContext := sdk.Context(user, fabsdk.WithOrg(orgName))

	resMgtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Printf("Failed to create new resource management client:%v\n", err)
		return resmgmt.InstantiateCCResponse{}, err
	}

	_, err = resMgtClient.InstallCC(resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg}, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Printf("InstallCC failed:%v\n", err)
		return resmgmt.InstantiateCCResponse{}, err
	}

	//ccPolicy := cauthdsl.SignedByMspAdmin(mspID)
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP", "Org2MSP"})
	req := resmgmt.InstantiateCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Args: ccArgs, Policy: ccPolicy}
	log.Printf("req: %+v\n", req)
	return resMgtClient.InstantiateCC(channelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
}

//GetOrgPeers ....
func GetOrgPeers(ctxProvider contextAPI.ClientProvider, org string) (peers []pfab.Peer, err error) {
	ctx, err := ctxProvider()
	if err != nil {
		return nil, err
	}
	orgPeers, err := ctx.EndpointConfig().PeersConfig(org)
	if err != nil {
		return nil, err
	}
	for _, orgpeer := range orgPeers {
		peer, err := ctx.InfraProvider().CreatePeerFromConfig(&pfab.NetworkPeer{PeerConfig: orgpeer})
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	return peers, nil
}
