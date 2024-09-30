package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// NewDockerClient initializes a new Docker client
func NewDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(context.Background())
	return cli, nil
}

// ListUnusedImages returns a list of unused Docker images
func ListUnusedImages(cli *client.Client) ([]image.Summary, error) {
	images, err := cli.ImageList(context.Background(), image.ListOptions{All: true})
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

// RemoveUnusedImages deletes unused Docker images
func RemoveUnusedImages(cli *client.Client, images []image.Summary) {
	opts := image.RemoveOptions{Force: true}

	for _, image := range images {
		_, err := cli.ImageRemove(context.Background(), image.ID, opts)
		if err != nil {
			log.Printf("Failed to remove image %s: %v", image.ID, err)
		} else {
			log.Printf("Successfully removed image %s", image.ID)
		}
	}
}

// PrintUnusedImages lists the images that would be removed (Dry Run)
func PrintUnusedImages(images []image.Summary) {
	if len(images) == 0 {
		fmt.Println("No unused images found.")
		return
	}

	fmt.Println("The following images would be removed:")
	for _, image := range images {
		fmt.Printf("ID: %s, Created: %d\n", image.ID, image.Created)
	}
}
