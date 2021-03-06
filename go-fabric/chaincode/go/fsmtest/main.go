package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/looplab/fsm"
)

func InitFSM(initStatus string) *fsm.FSM {
	f := fsm.NewFSM(
		initStatus,
		fsm.Events{
			{Name: "Submit", Src: []string{"Draft"}, Dst: "Submited"},
			{Name: "Approve", Src: []string{"Submited"}, Dst: "L1Approved"},
			{Name: "Reject", Src: []string{"Submited"}, Dst: "Reject"},
			{Name: "Approve", Src: []string{"L1Approved"}, Dst: "Complete"},
			{Name: "Reject", Src: []string{"L1Approved"}, Dst: "Reject"},
		},
		fsm.Callbacks{},
	)
	return f
}

type SimpleChaincode struct{}

func (s *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "Draft" {
		return s.Draft(stub, args)
	} else if function == "Submit" {
		return FsmEvent(stub, args, "Submit")
	} else if function == "Approve" {
		return FsmEvent(stub, args, "Approve")
	} else if function == "Reject" {
		return FsmEvent(stub, args, "Reject")
	}
	return shim.Error("Received unknown function invocation")
}

func (s *SimpleChaincode) Draft(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	formNumber := args[0]
	status := "Draft"
	stub.PutState(formNumber, []byte(status)) //初始化Draft状态的表单保存到StateDB
	return shim.Success([]byte(status))
}

func FsmEvent(stub shim.ChaincodeStubInterface, args []string, event string) pb.Response {
	formNumber := args[0]
	bstatus, err := stub.GetState(formNumber) //从StateDB中读取对应表单的状态
	if err != nil {
		return shim.Error("Query form status fail, form number:" + formNumber)
	}
	status := string(bstatus)
	fmt.Println("Form[" + formNumber + "] status:" + status)
	f := InitFSM(status) //初始化状态机，并设置当前状态为表单的状态
	err = f.Event(event) //触发状态机的事件
	if err != nil {
		return shim.Error("Current status is " + status + " does not support envent:" + event)
	}
	status = f.Current()
	fmt.Println("New status:" + status)
	stub.PutState(formNumber, []byte(status)) //更新表单的状态
	return shim.Success([]byte(status))       //返回新状态
}

func main() {
	if err := shim.Start(new(SimpleChaincode)); err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
