package test

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartNpmRegistryContainer(ctx context.Context) (*NpmRegistryContainer, error) {
	npmRegistry, err := setupNpmRegistry(ctx)
	if err != nil {
		return nil, err
	}

	/*resp, err := http.Get(npmRegistry.URI)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}*/

	return npmRegistry, nil
}

type NpmRegistryContainer struct {
	testcontainers.Container
	URI string
}

func setupNpmRegistry(ctx context.Context) (*NpmRegistryContainer, error) {
	rootPath := RootDir()
	req := testcontainers.ContainerRequest{
		Image:        "verdaccio/verdaccio",
		ExposedPorts: []string{"4873/tcp"},
		AutoRemove:   true,
		BindMounts: map[string]string{
			"/verdaccio/conf/config.yaml": rootPath + "/examples/verdaccio/config.yaml",
			"/verdaccio/storage/htpasswd": rootPath + "/examples/verdaccio/htpasswd",
			"/verdaccio/storage/data":     rootPath + "/examples/verdaccio/data",
		},
		WaitingFor: wait.ForLog("http://0.0.0.0:4873/"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "4873")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &NpmRegistryContainer{Container: container, URI: uri}, nil
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
