package main

import (
	"log"
	"sync"
	"time"

	"github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/integration/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

const (
	eventTimeout = time.Second * 30
)

var (
	org1Name          = "Org1"
	configFile        = "../../fixtures/config/config_test0.yaml"
	channelID         = "mychannel"
	channelConfigFile = "../../fixtures/fabric/v1.0.0/channel/mychannel.tx"
	chainCodePath     = "../chaincode"
	adminUser         = "Admin"
)
var eventCCArgs = [][]byte{[]byte("invoke"), []byte("SEVERE")}

func main() {
	chainCodeID := comm.GenerateRandomID()
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
	if err := comm.InstallAndInstantiateCC(c.SDK, fabsdk.WithUser(c.AdminUser), c.OrgID, c.ChannelID, chainCodeID, "github.com/events_cc", "v0", c.GoPath, nil); err != nil {
		log.Fatalf("InstallAndInstantiateCC return error: %v", err)
	}
	testReconnectEventHub(c)
	testFailedTx(c, chainCodeID)
	testFailedTxErrorCode(c, chainCodeID)
	testMultipleBlockEventCallbacks(c, chainCodeID)
}

func testFailedTx(c comm.Common, chainCodeID string) {
	fcn := "invoke"

	// Arguments for events CC
	var args [][]byte
	args = append(args, []byte("invoke"))
	args = append(args, []byte("SEVERE"))

	tpResponses1, prop1, err := comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, args, c.Targets[:1], nil)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal return error: %v", err)
	}

	tpResponses2, prop2, err := comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, args, c.Targets[:1], nil)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal return error: %v", err)
	}

	// Register tx1 and tx2 for commit/block event(s)
	// TransactionID提供Fabric交易提案的标识符。
	done1, fail1 := comm.RegisterTxEvent(prop1.TxnID, c.EventHub)
	defer c.EventHub.UnregisterTxEvent(prop1.TxnID)

	done2, fail2 := comm.RegisterTxEvent(prop2.TxnID, c.EventHub)
	defer c.EventHub.UnregisterTxEvent(prop2.TxnID)
	//监听事件
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorFailedTx(c, done1, fail1, done2, fail2)
	}()
	//测试无效事务：快速连续创建2个调用请求，修改相同的状态变量，导致一个调用无效
	_, err = comm.CreateAndSendTransaction(c.Transactor, prop1, tpResponses1)
	if err != nil {
		log.Fatalf("First invoke failed err: %v", err)
	}
	_, err = comm.CreateAndSendTransaction(c.Transactor, prop2, tpResponses2)
	if err != nil {
		log.Fatalf("Second invoke failed err: %v", err)
	}
	wg.Wait()
}

func testFailedTxErrorCode(c comm.Common, chainCodeID string) {
	fcn := "invoke"

	tpResponse1, prop1, err := comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, eventCCArgs, c.Targets[:1], nil)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal return error: %v", err)
	}

	tpResponse2, prop2, err := comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, eventCCArgs, c.Targets[:1], nil)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal return error: %v", err)
	}

	done := make(chan bool)
	fail := make(chan pb.TxValidationCode)

	c.EventHub.RegisterTxEvent(prop1.TxnID, func(txId fab.TransactionID, errorCode pb.TxValidationCode, err error) {
		if err != nil {
			fail <- errorCode
		} else {
			done <- true
		}
	})
	defer c.EventHub.UnregisterTxEvent(prop1.TxnID)

	done2 := make(chan bool)
	fail2 := make(chan pb.TxValidationCode)

	c.EventHub.RegisterTxEvent(prop2.TxnID, func(txId fab.TransactionID, errorCode pb.TxValidationCode, err error) {
		if err != nil {
			fail2 <- errorCode
		} else {
			done2 <- true
		}
	})
	defer c.EventHub.UnregisterTxEvent(prop2.TxnID)

	//monitoring of events
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorFailedTxErrorCode(c, done, fail, done2, fail2)
	}()
	//测试无效事务：快速连续创建2个调用请求，修改相同的状态变量，导致一个调用无效
	_, err = comm.CreateAndSendTransaction(c.Transactor, prop1, tpResponse1)
	if err != nil {
		log.Fatalf("First invoke failed err: %v", err)
	}
	_, err = comm.CreateAndSendTransaction(c.Transactor, prop2, tpResponse2)
	if err != nil {
		log.Fatalf("Second invoke failed err: %v", err)
	}
	wg.Wait()
}

func testMultipleBlockEventCallbacks(c comm.Common, chainCodeID string) {
	fcn := "invoke"
	//创建并注册将在块事件中调用的测试回调
	test := make(chan bool)
	c.EventHub.RegisterBlockEvent(func(block *common.Block) {
		log.Printf("Received test callback on block event")
		test <- true
	})
	//Transactor 交易者提供发送交易提案和交易的方法。
	tpResponses, prop, err := comm.CreateAndSendTransactionProposal(c.Transactor, chainCodeID, fcn, eventCCArgs, c.Targets[:1], nil)
	if err != nil {
		log.Fatalf("CreateAndSendTransactionProposal returned error: %v", err)
	}

	//注册tx提交/阻塞事件（s）
	done, fail := comm.RegisterTxEvent(prop.TxnID, c.EventHub)
	defer c.EventHub.UnregisterTxEvent(prop.TxnID)

	//up  monitoring of events
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorMultipleBlockEventCallbacks(c, done, fail, test)
	}()

	_, err = comm.CreateAndSendTransaction(c.Transactor, prop, tpResponses)
	if err != nil {
		log.Fatalf("CreateAndSendTransaction failed with error: %v", err)
	}
	wg.Wait()
}

func monitorFailedTx(c comm.Common, done1 chan bool, fail1 chan error, done2 chan bool, fail2 chan error) {
	rcvDone := false
	rcvFail := false
	timeout := time.After(eventTimeout)

Loop:
	for !rcvDone || !rcvFail {
		select {
		case <-done1:
			rcvDone = true
		case <-fail1:
			log.Fatalf("Received fail for first invoke")
		case <-done2:
			log.Fatalf("Received success for second invoke")
		case <-fail2:
			rcvFail = true
		case <-timeout:
			log.Println("Timeout: Didn't receive events")
			break Loop
		}
	}

	if !rcvDone || !rcvFail {
		log.Fatalf("Didn't receive events (done: %t; fail %t)", rcvDone, rcvFail)
	}
}

func monitorFailedTxErrorCode(c comm.Common, done chan bool, fail chan pb.TxValidationCode, done2 chan bool, fail2 chan pb.TxValidationCode) {
	rcvDone := false
	rcvFail := false
	timeout := time.After(eventTimeout)

Loop:
	for !rcvDone || !rcvFail {
		select {
		case <-done:
			rcvDone = true
		case <-fail:
			log.Fatalf("Received fail for first invoke")
		case <-done2:
			log.Fatalf("Received success for second invoke")
		case errorValidationCode := <-fail2:
			log.Printf("fail2 errorCode: %s\n", errorValidationCode)
			if errorValidationCode.String() != "MVCC_READ_CONFLICT" {
				log.Fatalf("Expected error code MVCC_READ_CONFLICT. Got %s", errorValidationCode.String())
			}
			rcvFail = true
		case <-timeout:
			log.Println("Timeout: Didn't receive events")
			break Loop
		}
	}

	if !rcvDone || !rcvFail {
		log.Fatalf("Didn't receive events (done: %t; fail %t)", rcvDone, rcvFail)
	}
}

func monitorMultipleBlockEventCallbacks(c comm.Common, done chan bool, fail chan error, test chan bool) {
	rcvTxDone := false
	rcvTxEvent := false
	timeout := time.After(eventTimeout)

Loop:
	for !rcvTxDone || !rcvTxEvent {
		select {
		case <-done:
			rcvTxDone = true
		case <-fail:
			log.Fatalf("Received tx failure")
		case <-test:
			rcvTxEvent = true
		case <-timeout:
			log.Printf("Timeout while waiting for events")
			break Loop
		}
	}

	if !rcvTxDone || !rcvTxEvent {
		log.Fatalf("Didn't receive events (tx event: %t; tx done %t)", rcvTxEvent, rcvTxDone)
	}
}

func testReconnectEventHub(c comm.Common) {
	// Test disconnect event hub
	if err := c.EventHub.Disconnect(); err != nil {
		log.Fatalf("Error disconnecting event hub: %s", err)
	}
	if c.EventHub.IsConnected() {
		log.Fatalf("Failed to disconnect event hub")
	}
	if err := c.EventHub.Connect(); err != nil {
		log.Fatalf("Failed to connect event hub")
	}
}
