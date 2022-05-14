package client

import (
	"context"
	"fmt"
	. "github.com/mohamed-gara/terraform-provider-npm/test"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNpmRegistryClient(t *testing.T) {
	ctx := context.Background()
	registry, registryErr := StartNpmRegistryContainer(ctx)
	if registryErr != nil {
		t.Fatal(registryErr)
	}
	defer registry.Terminate(ctx)

	client, clientErr := clientFor(registry)
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

func TestDownloadedArchiveCheckSum(t *testing.T) {
	packageName := "package-name"
	checksum := "b3268f732541120317f61a3227cdb2a86d63a897"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/package-name" {
			w.WriteHeader(http.StatusOK)

			const body = `
				{
				  "versions": {
					"2.0.0": {
					  "Dist": {
						"Shasum": "%s",
						"Tarball": "http://%s/my-package-npm-1.1.0.tgz"
					  }
					}
				  }
				}
			`
			finalBody := fmt.Sprintf(body, checksum, r.Host)
			w.Write([]byte(finalBody))
		}

		if r.URL.Path == "my-package-npm-1.1.0.tgz" {
			w.Write([]byte(`Unmatching body`))
		}
	}))
	defer server.Close()

	client, _ := authenticatedClientFor(server.URL)

	error := client.DownloadPackage(packageName, "2.0.0", ".")

	if error == nil {
		t.Fatalf("We should have error when downloading package with sha1 error.")
	}

	expectedErrorMsg := fmt.Sprintf("Downloaded file from %s/my-package-npm-1.1.0.tgz should have '%s' as checksum instead of 'da39a3ee5e6b4b0d3255bfef95601890afd80709'", server.URL, checksum)
	if error.Error() != expectedErrorMsg {
		t.Fatalf("Error should be '%s' instead of: '%s'.", expectedErrorMsg, error)
	}
}

func clientFor(registry *NpmRegistryContainer) (*NpmRegistry, error) {
	return authenticatedClientFor(registry.URI)
}

func authenticatedClientFor(url string) (*NpmRegistry, error) {
	user := "user"
	secret := "verdaccio"
	return NewNpmRegistry(&url, &user, &secret)
}
