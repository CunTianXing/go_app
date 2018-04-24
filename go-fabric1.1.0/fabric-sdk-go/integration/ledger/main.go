package main

import (
	"log"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
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

func main() {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatal(err)
	}
	defer sdk.Close()
	actionLedger(sdk)
}

func actionLedger(sdk *fabsdk.FabricSDK) {
	org1AdminChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))

	//Ledger client
	client, err := ledger.New(org1AdminChannelContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}

	ledgerInfo, err := client.QueryInfo()
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}
	log.Printf("ledger info : %+v\n", ledgerInfo)
	log.Printf("ledger ledgerInfo.BCI: %+v\n", ledgerInfo.BCI)
	log.Printf("ledger ledgerInfo.BCI.GetHeight: %+v\n", ledgerInfo.BCI.GetHeight())
	log.Printf("ledger ledgerInfo.BCI.GetPreviousBlockHash: %+v\n", string(ledgerInfo.BCI.PreviousBlockHash))
	log.Printf("ledger ledgerInfo.BCI.GetCurrentBlockHash: %+v\n", string(ledgerInfo.BCI.CurrentBlockHash))
	log.Printf("ledger ledgerInfo.Endorser: %+v\n", ledgerInfo.Endorser)
	log.Printf("ledger ledgerInfo.Status: %+v\n", ledgerInfo.Status)
	log.Println("=======================")
	configBackend, err := sdk.Config()
	if err != nil {
		log.Fatalf("failed to get config backend, error: %+v", err)
	}

	endpointConfig, err := fab.ConfigFromBackend(configBackend)
	if err != nil {
		log.Fatalf("failed to get endpoint config, error: %v", err)
	}
	log.Printf("endpointConfig :%+v\n", endpointConfig)
	log.Printf("endpointConfig  endpointConfig.CryptoConfigPath: %+v\n", endpointConfig.CryptoConfigPath())
	org2AdminClientContext := sdk.Context(fabsdk.WithUser(org2AdminUser), fabsdk.WithOrg(org2))

	peers, err := integration.GetOrgPeers(org2AdminClientContext, org2)
	if err != nil {
		log.Fatalf("get org2 peers failed: %+v\n", err)
	}
	log.Printf("get org2 peers: %+v\n", peers)

	ledgerInfoFromTarget, err := client.QueryInfo(ledger.WithTargets(peers[0]))
	if err != nil {
		log.Fatalf("QueryInfo return error: %v", err)
	}
	log.Printf("ledger info : %+v\n", ledgerInfoFromTarget)
	log.Printf("ledger ledgerInfo.BCI: %+v\n", ledgerInfoFromTarget.BCI)
	log.Printf("ledger ledgerInfo.BCI.GetHeight: %+v\n", ledgerInfoFromTarget.BCI.GetHeight())
	log.Printf("ledger ledgerInfo.BCI.GetPreviousBlockHash: %+v\n", string(ledgerInfoFromTarget.BCI.PreviousBlockHash))
	log.Printf("ledger ledgerInfo.BCI.GetCurrentBlockHash: %+v\n", string(ledgerInfoFromTarget.BCI.CurrentBlockHash))
	log.Printf("ledger ledgerInfo.Endorser: %+v\n", ledgerInfoFromTarget.Endorser)
	log.Printf("ledger ledgerInfo.Status: %+v\n", ledgerInfoFromTarget.Status)
	log.Println("================")

	block, err := client.QueryBlockByHash(ledgerInfo.BCI.CurrentBlockHash)
	if err != nil {
		log.Fatalf("QueryBlockByHash return error: %v", err)
	}
	log.Printf("QueryBlockByHash return :%+v\n", block)
	for index := 0; index < int(ledgerInfoFromTarget.BCI.Height); index++ {
		log.Println("=========", index, "=======")
		block, err = client.QueryBlock(0)
		if err != nil {
			log.Fatalf("QueryBlock return error: %v", err)
		}
		log.Printf("QueryBlock %d, GetPreviousHash :%s\n", index, string(block.GetHeader().GetPreviousHash()))
		log.Printf("QueryBlock %d, GetDataHash :%s\n", index, string(block.GetHeader().GetDataHash()))
		//log.Printf("QueryBlock %d, return :%+v\n", index, block)
	}

}
