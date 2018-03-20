package main

import (
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	cryptosuite "github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/sw"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/identitymgr"
	kvs "github.com/hyperledger/fabric-sdk-go/pkg/fab/keyvaluestore"
)

var (
	configFile       = "../../fixtures/config/config_test0.yaml"
	org1Name         = "Org1"
	org2Name         = "Org2"
	testFabricConfig core.Config
)

func main() {
	testFabricConfig, err := config.FromFile(configFile)()
	if err != nil {
		log.Fatalf("Failed InitConfig [%s]\n", err)
	}
	cryptoSuiteProvider, err := cryptosuite.GetSuiteByConfig(testFabricConfig)
	if err != nil {
		log.Fatalf("Failed getting cryptosuite from config : %s", err)
	}

	stateStore, err := kvs.New(&kvs.FileKeyValueStoreOptions{Path: testFabricConfig.CredentialStorePath()})
	if err != nil {
		log.Fatalf("CreateNewFileKeyValueStore failed: %v", err)
	}

	caClient, err := identitymgr.New(org2Name, stateStore, cryptoSuiteProvider, testFabricConfig)
	if err != nil {
		log.Fatalf("NewFabricCAClient return error: %v", err)
	}

	err = caClient.Enroll("admin", "adminpw")
	if err != nil {
		log.Fatalf("Enroll returned error: %v", err)
	}
}
