package main

import (
	"fmt"

	fsgConfig "github.com/hyperledger/fabric-sdk-go/pkg/config"
)

func main() {
	configImpl := fsgConfig.FromFile("config.yaml")
	fmt.Println(configImpl)
}
