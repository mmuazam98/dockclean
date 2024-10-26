package docker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

	const (
		tableline   = "----------------------------------------------------------------------------------------------------------------------------------------------------"
		tableformat = "%-35s %-12s %-30s %-15s %-30s\n"
	)

	// Print table header
	fmt.Println(tableline)
	fmt.Printf(tableformat, "ID", "Size", "Created (RFC3339)", "Status", "Labels")
	fmt.Println(tableline)

	// Iterate over each unused image and attempt removal
	for _, image := range images {
		// Remove the image
		_, err := d.CLI.ImageRemove(context.Background(), image.ID, opts)
		if err != nil {
			log.Printf("Failed to remove image %s: %v\n", image.ID, err)
		} else {
			// timestamp in RFC3339 format
			created := time.Unix(image.Created, 0).Format(time.RFC3339)
			truncatedDockerImageID := truncateDockerImageID(image.ID, 32)

			// Print image information in a table-like format
			fmt.Printf(tableformat,
				truncatedDockerImageID,
				formatSize(image.Size),
				created,
				"Removed",
				formatLabels(image.Labels),
			)
		}
	}

}

// Helper function to truncate strings with ellipsis
func truncateDockerImageID(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper function to converts bytes to human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Helper function used to format the docker image labels Key Value Pairs
func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "No Labels Found"
	}

	var labelStore []string
	for labelKey, labelValue := range labels {
		// Handling empty values
		if labelValue == "" {
			labelStore = append(labelStore, labelKey)
			continue
		}
		labelStore = append(labelStore, fmt.Sprintf("%s:%s", labelKey, labelValue))
	}

	return strings.Join(labelStore, ", ")
}

// Remove Images that exceed a specific size limit
func (d *DockerClient) RemoveExceedSizeLimit(sizeLimit int64, unit string) {

	var sizeLimitInBytes int64 = toBytes(sizeLimit, unit)

	images, err := d.ListUnusedImages()
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	if len(images) == 0 {
		log.Println("No unused images found")
		return
	}

	opts := image.RemoveOptions{Force: true}

	removedImagesCount := 0
	totalSizeCleaned := int64(0)

	for _, image := range images {

		// checking and removing images exceeding the threshold size
		if image.Size >= sizeLimitInBytes {
			_, err := d.CLI.ImageRemove(context.Background(), image.ID, opts)
			if err != nil {
				log.Printf("Failed to remove image %s: %v", image.ID, err)
			} else {
				log.Printf("Successfully removed image %s", image.ID)
			}

			totalSizeCleaned += image.Size
			removedImagesCount++
		}

	}

	if removedImagesCount > 0 {
		log.Printf("Summary: Removed %d images (Total space freed: %s)", removedImagesCount, formatSize(totalSizeCleaned))
	} else {
		fmt.Printf("No Unused Images are exceeding the limit %d %s", sizeLimit, strings.ToUpper(unit))
	}

}

// converts any value of any value to bytes
func toBytes(size int64, unit string) int64 {

	var multipliers map[string]int64 = map[string]int64{
		"b":  1,
		"kb": 1024,
		"mb": 1024 * 1024,
		"gb": 1024 * 1024 * 1024,
	}

	return size * multipliers[unit]

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
