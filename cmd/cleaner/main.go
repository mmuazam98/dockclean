package main

import (
	"flag"
	"log"

	"github.com/mmuazam98/dockclean/pkg/docker"
)

func parseFlag(dc docker.DockerClient) {

	dryRun := flag.Bool("dry-run", false, "List unused Docker images without deleting them")
	removeStopped := flag.Bool("remove-stopped", false, "Remove Images Associated with Stopped Containers")

	flag.Parse()

	switch {

	case *dryRun:
		dc.PrintUnusedImages()
	case *removeStopped:
		dc.CleanupStoppedContainerImages()
	default:
		dc.RemoveUnusedImages()

	}
}

func main() {

	var dc docker.DockerClient

	err := dc.InitDockerClient()
	if err != nil {
		log.Fatalf("Error initializing docker client : %v", err)
	}

	parseFlag(dc)
}
