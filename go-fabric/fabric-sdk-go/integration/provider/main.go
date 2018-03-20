package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/comm"
	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/params"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defsvc"

	selection "github.com/hyperledger/fabric-sdk-go/pkg/client/common/selection/dynamicselection"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
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

	chainCodeID := comm.GenerateRandomID()
	ccPath := "github.com/example_cc"
	if err := comm.InstallAndInstantiateCC(c.SDK, fabsdk.WithUser(c.AdminUser), c.OrgID, c.ChannelID, chainCodeID, ccPath, "v0", c.GoPath, params.ExampleCCInitArgs()); err != nil {
		log.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
	}
	//指定动态选择服务将使用的用户（检索chanincode策略信息）
	//此用户必须具有查询链接代码数据的lscc的权限
	//ChannelUser包含用于特定频道的用户（身份）信息
	mychannelUser := selection.ChannelUser{ChannelID: c.ChannelID, UserName: "User1", OrgName: c.OrgID}

	//使用动态选择为通道客户端创建SDK设置
	sdk, err := fabsdk.New(config.FromFile(c.ConfigFile),
		fabsdk.WithServicePkg(&DynamicSelectionProviderFactory{ChannelUsers: []selection.ChannelUser{mychannelUser}}))
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()

	chClient, err := sdk.NewClient(fabsdk.WithUser("User1")).Channel(c.ChannelID)
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}
	response, err := chClient.Query(channel.Request{ChaincodeID: chainCodeID, Fcn: "invoke", Args: params.ExampleCCQueryArgs()})
	if err != nil {
		log.Fatalf("Failed to query funds: %s", err)
	}
	value := response.Payload

	// Move funds
	response, err = chClient.Execute(channel.Request{ChaincodeID: chainCodeID, Fcn: "invoke", Args: params.ExampleCCTxArgs()})
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
	}

	// Verify move funds transaction result
	response, err = chClient.Query(channel.Request{ChaincodeID: chainCodeID, Fcn: "invoke", Args: params.ExampleCCQueryArgs()})
	if err != nil {
		log.Fatalf("Failed to query funds after transaction: %s", err)
	}
	valueInt, _ := strconv.Atoi(string(value))
	valueAfterInvokeInt, _ := strconv.Atoi(string(response.Payload))
	if valueInt+1 != valueAfterInvokeInt {
		log.Fatalf("Execute failed. Before: %s, after: %s", value, response.Payload)
	}
	fmt.Println("ok")
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
