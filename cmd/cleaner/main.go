package main

import (
	"flag"
	"log"

	"github.com/docker/docker/client"
	"github.com/mmuazam98/dockclean/pkg/docker"
)

func parseFlag(cli *client.Client) {
	dryRun := flag.Bool("dry-run", false, "List unused Docker images without deleting them")
	removeStopped := flag.Bool("remove-stopped", false, "Remove Images Associated with Stopped Containers")

	flag.Parse()

	if *dryRun {
		docker.PrintUnusedImages(cli)
	} else if *removeStopped {
		docker.CleanupStoppedContainerImages(cli)
	} else {
		docker.RemoveUnusedImages(cli)
	}
}

func main() {
	cli, err := initDockerClient()
	if err != nil {
		log.Fatalf("Failed to initialize Docker client: %v", err)
	}

	parseFlag(cli)
}

func initDockerClient() (*client.Client, error) {
	cli, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}
	return cli, nil
}
