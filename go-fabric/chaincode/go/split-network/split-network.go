package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Define Status codes for the response
const (
	OK    = 200
	ERROE = 500
)

// SmartContract struct
type SmartContract struct{}

// Init is called when the smart contract is instantiated
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SmartContract) update(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, expecting 3")
	}
	name := args[0]
	op := args[2]
	_, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Provided value was not a number")
	}
	if op != "+" && op != "-" {
		return shim.Error(fmt.Sprintf("Operator %s is unrecognized", op))
	}
	txid := APIstub.GetTxID()
	fmt.Println("GetTxID: ", txid)
	//GetTxID:  b10d58328fff8abb985da3fad32c30da92f362d2414b4bcc65218681b34548d3
	compositeIndexName := "varName~op~value~txID"
	attributes := []string{name, op, args[1], txid}
	fmt.Printf("CreateCompositeKey attributes: %#v\n", attributes)
	//CreateCompositeKey attributes: []string{"myvar", "+", "100", "b10d58328fff8abb985da3fad32c30da92f362d2414b4bcc65218681b34548d3"}
	compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, attributes)
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}
	fmt.Printf("CreateCompositeKey return: %#v\n", compositeKey)
	//CreateCompositeKey return: "\x00varName~op~value~txID\x00myvar\x00+\x00100\x00b10d58328fff8abb985da3fad32c30da92f362d2414b4bcc65218681b34548d3\x00"
	compositePutErr := APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}
	return shim.Success([]byte(fmt.Sprintf("Successfully added %s%s to %s", op, args[1], name)))
}

func (s *SmartContract) get(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	name := args[0]
	keys := []string{name}
	fmt.Printf("GetStateByPartialCompositeKey keys: %#v\n", keys)
	//GetStateByPartialCompositeKey keys: []string{"myvar"}
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~op~value~txID", keys)
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	fmt.Printf("GetStateByPartialCompositeKey return: %#v\n", deltaResultsIterator)
	//GetStateByPartialCompositeKey return: &shim.StateQueryIterator{CommonIterator:(*shim.CommonIterator)(0xc420159d70)}
	defer deltaResultsIterator.Close()
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}
	var finalVal float64
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}
		fmt.Printf("deltaResultsIterator.Next() return: %#v\n", responseRange)
		//deltaResultsIterator.Next() return: &queryresult.KV{Namespace:"bigdatacc", Key:"\x00varName~op~value~txID\x00myvar\x00+\x00100\x00b10d58328fff8abb985da3fad32c30da92f362d2414b4bcc65218681b34548d3\x00", Value:[]uint8{0x0}}
		fmt.Printf("responseRange.Key : %#v\n", responseRange.Key)
		//responseRange.Key : "\x00varName~op~value~txID\x00myvar\x00+\x00100\x00b10d58328fff8abb985da3fad32c30da92f362d2414b4bcc65218681b34548d3\x00"
		key, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}
		fmt.Printf("APIstub.SplitCompositeKey return: %#v, %#v\n", key, keyParts)
		//APIstub.SplitCompositeKey return: "varName~op~value~txID", []string{"myvar", "+", "100", "b10d58328fff8abb985da3fad32c30da92f362d2414b4bcc65218681b34548d3"}
		operation := keyParts[1]
		valueStr := keyParts[2]
		fmt.Printf("valueStr: %#v\n", valueStr)
		//valueStr: "100"
		value, convErr := strconv.ParseFloat(valueStr, 64)
		if convErr != nil {
			return shim.Error(convErr.Error())
		}
		switch operation {
		case "+":
			finalVal += value
		case "-":
			finalVal -= value
		default:
			return shim.Error(fmt.Sprintf("Unrecognized operation %s", operation))
		}
	}
	fmt.Printf("finalVal: %#v\n", finalVal)
	//finalVal: 100
	return shim.Success([]byte(strconv.FormatFloat(finalVal, 'f', -1, 64)))
}

func (s *SmartContract) delete(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	name := args[0]

	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~op~value~txID", []string{name})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve delta rows for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(fmt.Sprintf("Could not retrieve next delta row: %s", nextErr.Error()))
		}
		deltaRowDelErr := APIstub.DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return shim.Error(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}
	}
	fmt.Println("=========ok")
	return shim.Success([]byte(fmt.Sprintf("Deleted %s, %d rows removed", name, i)))
}

// Invoke routes invocations to the appropriate function in chaincode
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "update" {
		return s.update(APIstub, args)
	} else if function == "get" {
		return s.get(APIstub, args)
	} else if function == "delete" {
		return s.delete(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	if err := shim.Start(new(SmartContract)); err != nil {
		fmt.Printf("Error create new Start Contract: %s", err)
	}

}
