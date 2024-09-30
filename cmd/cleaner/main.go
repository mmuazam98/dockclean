package main

import (
	"flag"
	"log"

	"github.com/mmuazam98/dockclean/pkg/docker"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "List unused Docker images without deleting them")
	flag.Parse()

	cli, err := docker.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed to initialize Docker client: %v", err)
	}

	// List and clean unused Docker images
	unusedImages, err := docker.ListUnusedImages(cli)
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	if *dryRun {
		docker.PrintUnusedImages(unusedImages)
	} else {
		docker.RemoveUnusedImages(cli, unusedImages)
	}
}
