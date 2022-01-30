package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/codeclysm/extract"
)

type NpmRegistryClient struct {
	httpClient  *http.Client
	url         string
	credentials credentials
}

type credentials struct {
	username string
	password string
}

type PackageMetadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	Dist struct {
		Shasum  string `json:"shasum"`
		Tarball string `json:"tarball"`
	} `json:"dist"`
}

func NewNpmRegistryClient(host, username, password *string) (*NpmRegistryClient, error) {
	client := NpmRegistryClient{
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

func (r *NpmRegistryClient) getPackageVersion(name, version string) (*PackageMetadata, error) {
	metadataUrl := r.url + "/" + name + "/" + version
	req, requestErr := http.NewRequest(http.MethodGet, metadataUrl, http.NoBody)
	if requestErr != nil {
		return nil, fmt.Errorf("failed to get packet information: %v", requestErr)
	}

	body, err := r.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get packet information: %v", err)
	}

	var pack *PackageMetadata
	decoder := json.NewDecoder(bytes.NewReader(body))
	decodeErr := decoder.Decode(&pack)
	if decodeErr != nil {
		return nil, fmt.Errorf("invalid response: %v", err)
	}

	return pack, nil
}

func (c NpmRegistryClient) DownloadPackageToDestination(name, version, destinationDirectory string) error {
	// Create the file
	/*fileName := fmt.Sprintf("%s-%s.tgz", name, version)
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	*/

	// remove old data if exist
	if _, err := os.Stat(destinationDirectory); err == nil {
		// zap.S().Debugf("remove old.")
		path, _ := filepath.Abs(destinationDirectory + "/package")
		log.Printf("Destination is: %s", path)
		err = os.RemoveAll(destinationDirectory + "/package")
		if err != nil {
			return fmt.Errorf("failed to remove old data: %v", err)
		}
	}

	metadata, getMetaDataErr := c.getPackageVersion(name, version)
	if getMetaDataErr != nil {
		return getMetaDataErr
	}

	// Get the data
	req, requestErr := http.NewRequest(http.MethodGet, metadata.Dist.Tarball, http.NoBody)
	if requestErr != nil {
		return fmt.Errorf("failed to download package tarball: %v", requestErr)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("failed to download package tarball: %v", err)
	}

	// extract archive
	err = extract.Archive(context.Background(), bytes.NewReader(body), destinationDirectory, func(name string) string {
		return name
	})
	if err != nil {
		return fmt.Errorf("failed to extract tar to destinationDirectory: %s", err)
	}

	return nil
}

func (c *NpmRegistryClient) doRequest(req *http.Request) ([]byte, error) {
	// req.Header.Set("Authorization", c.token)
	req.SetBasicAuth(c.credentials.username, c.credentials.password)

	res, err := c.httpClient.Do(req)
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

/*func tokenOf(username *string, password *string) string {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(*password))
	return "Basic " + *username + ":" + encodedPassword
}
*/
