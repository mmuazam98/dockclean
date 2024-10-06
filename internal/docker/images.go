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

// VerboseModeCleanup gives more details while doing the cleanup of unused images
func (d *DockerClient) VerboseModeCleanup() {

	images, err := d.ListUnusedImages()
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	opts := image.RemoveOptions{Force: true}

	if len(images) == 0 {
		log.Println("No unused images found")
		return
	}

	log.Printf("Found %d unused images. Starting removal in verbose mode...\n", len(images))

	// Print table header
	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("%-15s %-12s %-15s %s\n", "ID", "Size (bytes)", "Created (Unix)", "Labels")
	fmt.Println("---------------------------------------------------------------")

	// Iterate over each unused image and attempt removal
	for _, image := range images {
		// Remove the image
		_, err := d.CLI.ImageRemove(context.Background(), image.ID, opts)
		if err != nil {
			log.Printf("Failed to remove image %s: %v\n", image.ID, err)
		} else {
			// Print image information in a table-like format
			fmt.Printf("%-15s %-12d %-15d ", image.ID[:12], image.Size, image.Created)

			// Display labels
			if len(image.Labels) > 0 {
				labelStr := ""
				for key, value := range image.Labels {
					labelStr += fmt.Sprintf("%s:%s, ", key, value)
				}
				fmt.Println(labelStr[:len(labelStr)-2]) // Remove the last comma
			} else {
				fmt.Println("No labels")
			}
		}
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
