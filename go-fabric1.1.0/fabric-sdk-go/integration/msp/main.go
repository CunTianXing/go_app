package main

import (
	"fmt"
	"log"

	"github.com/CunTianXing/go_app/go-fabric1.1.0/fabric-sdk-go/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	IdentifyTypeUser = "User"
	configFile       = "../../fixtures/config/config_test.yaml"
)

func main() {
	configProvider := config.FromFile(configFile)
	//instantiate the sdk
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatal(err)
	}
	integration.CleanupUserData(sdk)
	defer integration.CleanupUserData(sdk)

	ctxProvider := sdk.Context()
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		log.Fatalf("failed to create CA client: %v", err)
	}
	registrarEnrollID, registrarEnrollSecret := getRegistrarEnrollmentCredentials(ctxProvider)
	fmt.Println(registrarEnrollSecret)
	fmt.Println(registrarEnrollID)
	err = mspClient.Enroll(registrarEnrollID, msp.WithSecret(registrarEnrollSecret))
	if err != nil {
		log.Fatalf("Enroll failed: %v", err)
	}

	// Generate a random user name
	username := integration.GenerateRandomID()

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        username,
		Type:        IdentifyTypeUser,
		Affiliation: "org2",
	})
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}
	// Enroll the new user
	err = mspClient.Enroll(username, msp.WithSecret(enrollmentSecret))
	if err != nil {
		log.Fatalf("Enroll failed: %v", err)
	}
	// Get the new user's signing identity
	_, err = mspClient.GetSigningIdentity(username)
	if err != nil {
		log.Fatalf("GetSigningIdentity failed: %v", err)
	}
}

func getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {
	ctx, err := ctxProvider()
	if err != nil {
		log.Fatalf("failed to get context: %v", err)
	}

	clientConfig, err := ctx.IdentityConfig().Client()
	if err != nil {
		log.Fatalf("config.Client() failed: %v", err)
	}

	myOrg := clientConfig.Organization
	fmt.Println(myOrg)
	caConfig, err := ctx.IdentityConfig().CAConfig(myOrg)
	if err != nil {
		log.Fatalf("CAConfig failed: %v", err)
	}
	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
}
