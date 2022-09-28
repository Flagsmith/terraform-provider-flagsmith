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

// Create provider factories - to be used by resource tests
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"flagsmith": providerserver.NewProtocol6WithError(flagsmith.New("")()),
}

func testAccPreCheck(t *testing.T) {
	mustHaveEnv(t, "FLAGSMITH_MASTER_API_KEY")
	mustHaveEnv(t, "FLAGSMITH_FEATURE_NAME")
	mustHaveEnv(t, "FLAGSMITH_ENVIRONMENT_KEY")
	mustHaveEnv(t, "FLAGSMITH_ENVIRONMENT_ID")
	mustHaveEnv(t, "FLAGSMITH_FEATURE_ID")
	mustHaveEnv(t, "FLAGSMITH_PROJECT_UUID")
	mustHaveEnv(t, "FLAGSMITH_FEATURE_UUID")
}

func mustHaveEnv(t *testing.T, name string) {
	if os.Getenv(name) == "" {
		t.Fatalf("%s environment variable must be set for acceptance tests", name)
	}
}

func masterAPIKey() string {
	return os.Getenv("FLAGSMITH_MASTER_API_KEY")
}
func featureName() string {
	return os.Getenv("FLAGSMITH_FEATURE_NAME")
}
func featureUUID() string {
	return os.Getenv("FLAGSMITH_FEATURE_UUID")
}
func projectUUID() string {
	return os.Getenv("FLAGSMITH_PROJECT_UUID")
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
	baseAPIURL := os.Getenv("FLAGSMITH_BASE_API_URL")
	if baseAPIURL == "" {
		baseAPIURL = "https://api.flagsmith.com/api/v1"
	}

	if tc == nil {
		tc = flagsmithapi.NewClient(masterAPIKey(), baseAPIURL)
	}

	return tc
}
