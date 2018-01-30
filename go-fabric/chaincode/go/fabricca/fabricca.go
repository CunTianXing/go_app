package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct{}

func (s *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	if function == "cert" {
		return s.testCertificate(stub, args)
	}
	return shim.Error("Received unknown function invocation")
}

func (s *SimpleChaincode) testCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----")
	if certStart == -1 {
		fmt.Println("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Println("Could not decode the PEM structure")
	}
	fmt.Println(string(certText))
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Println("ParseCertificate failed")
	}
	fmt.Println(cert)
	uname := cert.Subject.CommonName
	fmt.Println("Name:" + uname)
	return shim.Success([]byte("Called testCertificate " + uname))
}

func main() {
	if err := shim.Start(new(SimpleChaincode)); err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
