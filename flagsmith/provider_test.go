package flagsmith_test

import (
	"github.com/Flagsmith/flagsmith-go-api-client"
	"github.com/Flagsmith/terraform-provider-flagsmith/flagsmith"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"strconv"
	"testing"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"flagsmith": providerserver.NewProtocol6WithError(flagsmith.New("")()),
}

func mustHaveEnv(t *testing.T, name string) {
	if os.Getenv(name) == "" {
		t.Fatalf("%s environment variable must be set for acceptance tests", name)
	}
}
func testAccPreCheck(t *testing.T) {
	mustHaveEnv(t, "FLAGSMITH_MASTER_API_KEY")
	mustHaveEnv(t, "FLAGSMITH_FEATURE_NAME")
	mustHaveEnv(t, "FLAGSMITH_ENVIRONMENT_KEY")
	mustHaveEnv(t, "FLAGSMITH_ENVIRONMENT_ID")
	mustHaveEnv(t, "FLAGSMITH_FEATURE_ID")
}

func masterAPIKey() string {
	return os.Getenv("FLAGSMITH_MASTER_API_KEY")
}
func featureName() string {
	return os.Getenv("FLAGSMITH_FEATURE_NAME")
}
func environmentKey() string {
	return os.Getenv("FLAGSMITH_ENVIRONMENT_KEY")
}
func environmentID() int {
	v, err := strconv.Atoi(os.Getenv("FLAGSMITH_ENVIRONMENT_ID"))
	if err != nil {
		panic(err)

	}
	return v
}
func featureID() int {
	v, err := strconv.Atoi(os.Getenv("FLAGSMITH_FEATURE_ID"))
	if err != nil {
		panic(err)
	}
	return v
}

var tc *flagsmithapi.Client

func testClient() *flagsmithapi.Client {
	if tc == nil {
		tc = flagsmithapi.NewClient(masterAPIKey(), "http://localhost:8000/api/v1")
	}
	return tc
}
