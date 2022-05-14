package client

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codeclysm/extract"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type NpmRegistry struct {
	httpClient  *http.Client
	url         string
	credentials credentials
}

type credentials struct {
	username string
	password string
}

func NewNpmRegistry(host, username, password *string) (*NpmRegistry, error) {
	client := NpmRegistry{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: *host,
		credentials: credentials{
			*username,
			*password,
		},
	}
	return &client, nil
}

func (registry *NpmRegistry) DownloadPackage(name, version, destinationDirectory string) error {

	cleanUpErr := cleanUpExisting(destinationDirectory)
	if cleanUpErr != nil {
		return cleanUpErr
	}

	metadata, getMetaDataErr := registry.fetchPackageVersionMetadata(name, version)
	if getMetaDataErr != nil {
		return getMetaDataErr
	}

	body, err := registry.downloadTarball(metadata)
	if err != nil {
		return err
	}

	shaError := verifySha1(body, metadata)
	if shaError != nil {
		return shaError
	}

	err = extract.Archive(context.Background(), bytes.NewReader(body), destinationDirectory, func(name string) string {
		return name
	})
	if err != nil {
		return fmt.Errorf("failed to extract tar to destinationDirectory: %s", err)
	}

	return nil
}

func (registry *NpmRegistry) fetchPackageVersionMetadata(name, version string) (*PackageVersionMetadata, error) {
	metadataUrl := registry.url + "/" + name
	req, requestErr := http.NewRequest(http.MethodGet, metadataUrl, http.NoBody)
	if requestErr != nil {
		return nil, fmt.Errorf("failed to fetch package metadata: %v", requestErr)
	}

	body, err := registry.sendAuthenticatedRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get packet information: %v", err)
	}

	var packageMetadata *PackageMetadata
	decoder := json.NewDecoder(bytes.NewReader(body))
	decodeErr := decoder.Decode(&packageMetadata)
	if decodeErr != nil {
		return nil, fmt.Errorf("invalid response: %v", decodeErr)
	}

	versionMetadata, exists := packageMetadata.Versions[version]
	if !exists {
		return nil, fmt.Errorf("package version %s@%v not found", name, version)
	}

	return versionMetadata, nil
}

func verifySha1(requestBody []byte, metadata *PackageVersionMetadata) error {
	h := sha1.New()
	h.Write(requestBody)

	actualSha1 := hex.EncodeToString(h.Sum(nil))
	expectedSha1 := metadata.Dist.Shasum

	if expectedSha1 != actualSha1 {
		return errors.New(fmt.Sprintf("Downloaded file from %s should have '%s' as checksum instead of '%s'", metadata.Dist.Tarball, metadata.Dist.Shasum, actualSha1))
	}
	return nil
}

func (registry *NpmRegistry) downloadTarball(metadata *PackageVersionMetadata) ([]byte, error) {
	req, requestErr := http.NewRequest(http.MethodGet, metadata.Dist.Tarball, http.NoBody)
	if requestErr != nil {
		return nil, fmt.Errorf("failed to download package tarball: %v", requestErr)
	}

	body, err := registry.sendAuthenticatedRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download package tarball: %v", err)
	}
	return body, nil
}

func (registry *NpmRegistry) sendAuthenticatedRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", tokenOf(&registry.credentials.username, &registry.credentials.password))

	res, err := registry.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func cleanUpExisting(destinationDirectory string) error {
	if _, err := os.Stat(destinationDirectory); err == nil {
		path, _ := filepath.Abs(destinationDirectory + "/package")
		log.Printf("Destination is: %s", path)
		err = os.RemoveAll(destinationDirectory + "/package")
		if err != nil {
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}
	return nil
}

func tokenOf(username *string, password *string) string {
	value := *username + ":" + *password
	token := base64.StdEncoding.EncodeToString([]byte(value))
	return "Basic " + token
}
