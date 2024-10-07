package main

import (
	"log"

	"github.com/mmuazam98/dockclean/internal/docker"
	"github.com/mmuazam98/dockclean/pkg/utils"
)

func main() {

	var dc docker.DockerClient

	err := dc.InitDockerClient()
	if err != nil {
		log.Fatalf("Error initializing docker client : %v", err)
	}

	f := utils.ParseFlags()

	switch {
	case f.DryRun:
		dc.PrintUnusedImages()
	case f.RemoveStopped:
		dc.CleanupStoppedContainerImages()
	default:
		dc.RemoveUnusedImages()
	}

}
