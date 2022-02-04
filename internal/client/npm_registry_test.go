package client

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	. "github.com/mohamed-gara/terraform-provider-npm/test"
)

func TestNpmRegistryClient(t *testing.T) {
	ctx := context.Background()
	registry, registryErr := StartNpmRegistryContainer(ctx)
	if registryErr != nil {
		t.Fatal(registryErr)
	}
	defer registry.Terminate(ctx)

	client, clientErr := AuthenticatedClientFor(registry)
	if clientErr != nil {
		t.Fatal(clientErr)
	}

	if client.url != registry.URI {
		t.Fatal("Client url should be the same as registry url.")
	}

	//TODO add file to package directory

	downloadErr := client.DownloadPackage("my-package-npm", "1.1.0", ".")
	if downloadErr != nil {
		t.Fatal(downloadErr)
	}

	var fileWalkingErr = filepath.Walk("./package", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal("Package destination directory not found.")
		}
		if !info.IsDir() {
			if info.Name() != "file1.txt" && info.Name() != "file2.txt" && info.Name() != "package.json" {
				t.Fatal("File is not expected: " + info.Name())
			}
		}
		return nil
	})
	if fileWalkingErr != nil {
		log.Fatal(fileWalkingErr)
	}
}

func AuthenticatedClientFor(registry *NpmRegistryContainer) (*NpmRegistry, error) {
	secret := "verdaccio"
	return NewNpmRegistry(&registry.URI, &secret, &secret)
}
