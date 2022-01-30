package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type providerFactory = func() (*schema.Provider, error)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = func(registryUrl string) map[string]providerFactory {
	return map[string]providerFactory{
		"npm": func() (*schema.Provider, error) {
			provider := New("dev")()
			// provider.Configure(context.Context(), terraform.NewResourceConfigRaw(nil))
			return provider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	if os.Getenv("TERRAFORM_PROVIDER_NPM_USERNAME") == "" {
		t.Fatal("TERRAFORM_PROVIDER_NPM_USERNAME must be set for acceptance tests")
	}

	if os.Getenv("TERRAFORM_PROVIDER_NPM_PASSWORD") == "" {
		t.Fatal("TERRAFORM_PROVIDER_NPM_PASSWORD must be set for acceptance tests")
	}

	//err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
	//if err != nil {
	//:t.Fatal(err)
	//}
}
