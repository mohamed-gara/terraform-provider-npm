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
	os.Setenv("TERRAFORM_PROVIDER_NPM_USERNAME", "verdaccio")
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.npm_package.my_pkg", "name", regexp.MustCompile("my-package-npm")),
					resource.TestMatchResourceAttr("data.npm_package.my_pkg", "version", regexp.MustCompile("1.1.0")),
					resource.TestMatchResourceAttr("data.npm_package.my_pkg", "files.0", regexp.MustCompile(folderAbsolutePath+"/file1.txt")),
					resource.TestMatchResourceAttr("data.npm_package.my_pkg", "files.1", regexp.MustCompile(folderAbsolutePath+"/file2.txt")),
					resource.TestMatchResourceAttr("data.npm_package.my_pkg", "files.2", regexp.MustCompile(folderAbsolutePath+"/package.json")),
				),
			},
		},
	})
}
