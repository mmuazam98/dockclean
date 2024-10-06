package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

// CleanupStoppedContainerImages removes images associated with stopped containers
func (d *DockerClient) CleanupStoppedContainerImages() {

	const (
		ImageStateUnreferenced = -1
		ImageStateExited       = 0
		ImageStateRunning      = 1
	)

	// List all images
	images, err := d.CLI.ImageList(context.Background(), image.ListOptions{})
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
	containers, err := d.CLI.ContainerList(context.Background(), container.ListOptions{All: true})
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
			_, err := d.CLI.ImageRemove(context.Background(), imageId, opts)
			if err != nil {
				log.Printf("Failed to remove image %s: %v", imageId, err)
			} else {
				log.Printf("Successfully removed image %s", imageId)
				// Remove associated containers only if the image was successfully deleted
				for containerID, imageID := range containersToRemove {
					if imageID == imageId {
						if err := d.CLI.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true}); err != nil {
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
