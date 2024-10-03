package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
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

// CleanupStoppedContainerImages removes images associated with stopped containers
func CleanupStoppedContainerImages(cli *client.Client) {

	const (
		ImageStateUnreferenced = -1
		ImageStateExited       = 0
		ImageStateRunning      = 1
	)

	// List all images
	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	imagesForCleanup := make(map[string]int8)

	// -1 : not referenced
	// 0  : exited
	// 1  : other state
	for _, image := range images {
		imagesForCleanup[image.ID] = ImageStateUnreferenced // Marking all images as -1, currently not referenced
	}

	// List all containers
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	// Create a map to hold containers associated with images to be deleted
	containersToRemove := make(map[string]string) // container ID -> image ID

	// Unreferenced containers are not affected so can identify using -1
	// If all are exited, the state will update once to 0 and won't be affected any later
	// Even if one of the container is not exited, it will be updated to one
	for _, container := range containers {
		containerState, containerImageID := container.State, container.ImageID

		if containerState == "exited" && imagesForCleanup[containerImageID] != ImageStateRunning {

			// Mark image as exited if no running containers are using it, otherwise mark it as running
			imagesForCleanup[containerImageID] = ImageStateExited
			containersToRemove[container.ID] = containerImageID // Track the container for later removal
		} else {
			imagesForCleanup[containerImageID] = ImageStateRunning
		}
	}

	// Remove all the images for those whose all of the associated containers are stopped, and not just one
	opts := image.RemoveOptions{Force: true}

	for imageId, removalState := range imagesForCleanup {
		if removalState == ImageStateExited {
			_, err := cli.ImageRemove(context.Background(), imageId, opts)
			if err != nil {
				log.Printf("Failed to remove image %s: %v", imageId, err)
			} else {
				log.Printf("Successfully removed image %s", imageId)
				// Remove associated containers only if the image was successfully deleted
				for containerID, imageID := range containersToRemove {
					if imageID == imageId {
						if err := cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true}); err != nil {
							log.Printf("Failed to remove stopped container %s: %v", containerID, err)
						} else {
							log.Printf("Successfully removed stopped container %s", containerID)
						}
					}
				}
			}
		}
	}

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
func RemoveUnusedImages(cli *client.Client) {

	images, err := ListUnusedImages(cli)
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

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
func PrintUnusedImages(cli *client.Client) {

	images, err := ListUnusedImages(cli)
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
