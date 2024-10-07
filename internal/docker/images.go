package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/image"
)

// ListUnusedImages returns a list of unused Docker images
func (d *DockerClient) ListUnusedImages() ([]image.Summary, error) {

	images, err := d.CLI.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var unusedImages []image.Summary
	for _, image := range images {
		// Filter images without RepoTags (untagged or unused images)
		if len(image.RepoTags) == 0 {
			unusedImages = append(unusedImages, image)
		}
	}

	if len(unusedImages) > 0 {
		fmt.Printf("Found %d unused images\n", len(unusedImages))
	}
	return unusedImages, nil
}

// PrintUnusedImages lists the images that would be removed (Dry Run)
func (d *DockerClient) PrintUnusedImages() {

	images, err := d.ListUnusedImages()
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	if len(images) == 0 {
		fmt.Println("No unused images found.")
		return
	}

	fmt.Println("The following images would be removed:")
	for _, image := range images {
		fmt.Printf("ID: %s, Created: %d\n", image.ID, image.Created)
	}
}

// RemoveUnusedImages deletes unused Docker images
func (d *DockerClient) RemoveUnusedImages() {

	images, err := d.ListUnusedImages()
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	opts := image.RemoveOptions{Force: true}

	for _, image := range images {
		_, err := d.CLI.ImageRemove(context.Background(), image.ID, opts)
		if err != nil {
			log.Printf("Failed to remove image %s: %v", image.ID, err)
		} else {
			log.Printf("Successfully removed image %s", image.ID)
		}
	}
}
