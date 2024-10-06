package docker

import (
	"context"

	"github.com/docker/docker/client"
)

type DockerClient struct {
	CLI *client.Client
}

// InitDockerClient initializes a new Docker client
func (d *DockerClient) InitDockerClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil
	}
	cli.NegotiateAPIVersion(context.Background())
	d.CLI = cli
	return nil
}
