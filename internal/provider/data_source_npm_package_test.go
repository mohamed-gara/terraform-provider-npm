package provider

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/mohamed-gara/terraform-provider-npm/test"
)

func TestAccDataSourceNpmPackage(t *testing.T) {
	ctx := context.Background()
	registry, registryErr := StartNpmRegistryContainer(ctx)
	if registryErr != nil {
		t.Fatal(registryErr)
	}
	defer registry.Terminate(ctx)

	os.Setenv("TERRAFORM_PROVIDER_NPM_URL", registry.URI)
	os.Setenv("TERRAFORM_PROVIDER_NPM_USERNAME", "user")
	os.Setenv("TERRAFORM_PROVIDER_NPM_PASSWORD", "verdaccio")
	defer os.Unsetenv("TERRAFORM_PROVIDER_NPM_URL")
	defer os.Unsetenv("TERRAFORM_PROVIDER_NPM_USERNAME")
	defer os.Unsetenv("TERRAFORM_PROVIDER_NPM_PASSWORD")

	folderAbsolutePath, pathErr := filepath.Abs("./package")
	if pathErr != nil {
		t.Fatal(pathErr)
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories(registry.URI),
		Steps: []resource.TestStep{
			{
				Config: `
					data "npm_package" "my_pkg" {
				  		name    = "my-package-npm"
				  		version = "1.1.0"
					}
				`,
				Check: checkAttributes(map[string]string{
					"name":                  "my-package-npm",
					"version":               "1.1.0",
					"files.0.absolute_path": folderAbsolutePath + "/file1.txt",
					"files.0.mime_type":     "text/plain",
					"files.1.absolute_path": folderAbsolutePath + "/file2.txt",
					"files.1.mime_type":     "text/plain",
					"files.2.absolute_path": folderAbsolutePath + "/package.json",
					"files.2.mime_type":     "application/json",
				}),
			},
		},
	})
}

func checkAttributes(attributes map[string]string) resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0, len(attributes))
	for k, v := range attributes {
		checks = append(checks, checkAttribute(k, v))
	}
	return resource.ComposeTestCheckFunc(checks...)
}

func checkAttribute(attribute, valueRegex string) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr("data.npm_package.my_pkg", attribute, regexp.MustCompile(valueRegex))
}
